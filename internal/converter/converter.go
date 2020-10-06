package converter

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/sirupsen/logrus"
)

// Converter is everything you need to convert protos to JSON Schemas:
type Converter struct {
	// Allow NULL values for all properties. By default, JSON Schemas will reject NULL values.
	AllowNullValues bool
	// Allow additional properties. JSON Schemas will allow extra parameters, that are not specified in the schema.
	AllowAdditionalProperties bool
	// If the parameter is not set (default) the JSON Schema will allow both string and integers for 64 bit integers.
	// If it is set only integers are allowed.
	// The canonical JSON encoding of Proto3 converts int64, fixed64, uint6 to JSON strings.
	// When decoding JSON to ProtoBuf both numbers and strings are accepted.
	DisallowBigIntsAsStrings bool
	// Allow both enum names and integer values.
	AllowNumericEnumValues bool
	// Path to an OpenAPI file that will be merged with the generated schemas.
	OpenApiFile string
	// Create a single file instead of multiple files.
	// When JSON Schema mode is enabled (parameter open_api is not set) a single JSON Schema will be generated with the
	// given filename.
	// When OpenAPI mode is enabled (parameter open_api is set) this parameter is implicitly enabled and the
	// default filename is "openapi.json".
	outputFile string
	// If the parameter is set the field names from the ProtoBuf definition are used in the JSON Schema.
	// If the parameter is not set message field names are mapped to lowerCamelCase and become JSON object keys.
	UseProtoFieldNames bool
	logger             *logrus.Logger
	sourceInfo         *sourceCodeInfo
	generatorPlan      *generatorPlan
}

// New returns a configured *Converter:
func New(logger *logrus.Logger) *Converter {
	return &Converter{
		logger:        logger,
		generatorPlan: NewGeneratorPlan(),
	}
}

// ConvertFrom tells the convert to work on the given input:
func (c *Converter) ConvertFrom(rd io.Reader) (*plugin.CodeGeneratorResponse, error) {
	c.logger.Debug("Reading code generation request")
	input, err := ioutil.ReadAll(rd)
	if err != nil {
		c.logger.WithError(err).Error("Failed to read request")
		return nil, err
	}

	req := &plugin.CodeGeneratorRequest{}
	err = proto.Unmarshal(input, req)
	if err != nil {
		c.logger.WithError(err).Error("Can't unmarshal input")
		return nil, err
	}

	c.parseGeneratorParameters(req.GetParameter())

	c.logger.Debug("Converting input")
	return c.convert(req)
}

func (c *Converter) parseGeneratorParameters(parameters string) {
	for _, parameter := range strings.Split(parameters, ",") {
		switch {
		case parameter == "": // ignore empty parameter
		case parameter == "allow_null_values":
			c.AllowNullValues = true
		case parameter == "allow_additional_properties":
			c.AllowAdditionalProperties = true
		case parameter == "debug":
			c.logger.SetLevel(logrus.DebugLevel)
		case parameter == "disallow_bigints_as_strings":
			c.DisallowBigIntsAsStrings = true
		case parameter == "allow_numeric_enum_values":
			c.AllowNumericEnumValues = true
		case strings.HasPrefix(parameter, "out_file"):
			paramSplit := strings.Split(parameter, "=")
			if len(paramSplit) != 2 || len(paramSplit[1]) == 0 {
				c.logger.WithField("parameter", parameter).
					Warn("Invalid parameter. Expecting out_file=<filename>")
			} else {
				c.outputFile = paramSplit[1]
			}
		case strings.HasPrefix(parameter, "open_api_template"):
			paramSplit := strings.Split(parameter, "=")
			if len(paramSplit) != 2 || len(paramSplit[1]) == 0 {
				c.logger.WithField("parameter", parameter).
					Warn("Invalid parameter. Expecting open_api=<filename>")
			} else {
				c.OpenApiFile = paramSplit[1]
			}
		case parameter == "proto_fieldnames":
			c.UseProtoFieldNames = true
		default:
			c.logger.WithField("parameter", parameter).Warn("Unknown parameter")
		}
	}

	// If an OpenAPI document should be generated single file mode is always enabled. The default filename is
	// openapi.json.
	if c.outputFile == "" {
		c.outputFile = "openapi.json"
	}
}

func (c *Converter) convert(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	c.sourceInfo = newSourceCodeInfo(req.GetProtoFile())
	res := &plugin.CodeGeneratorResponse{}

	err := c.calcGeneratorPlan(req)
	if err != nil {
		return nil, err
	}
	converted, err := c.convertFile(c.outputFile)
	if err != nil {
		c.logger.WithField("file", c.outputFile).WithField("error", err).
			Error("Failed to create JSON schema file")
		res.Error = proto.String(fmt.Sprintf("Failed to create %s: %v", c.outputFile, err))
		return res, err
	}
	res.File = append(res.File, converted)
	return res, nil
}

func (c *Converter) convertFile(openapiFileName string) (*plugin.CodeGeneratorResponse_File, error) {
	typeInfos := c.generatorPlan.GetAll()

	openApi, err := c.createOpenApiDocument()
	if err != nil {
		return nil, err
	}

	c.logger.WithField("file", openapiFileName).Debug("Creating OpenAPI file ...")

	for _, typeInfo := range typeInfos {
		jsonType, err := c.convertProtoType(typeInfo)
		if err != nil {
			return nil, err
		}

		typeId := typeInfo.GetProtoFQTypeName()
		if openApi.Components.Schemas == nil {
			openApi.Components.Schemas = Schemas{}
		}
		openApi.Components.Schemas[typeId] = jsonType
	}

	// Marshal the result into a JSON-Schema or OpenAPI document:
	var jsonSchemaJSON []byte = nil
	jsonBytes, err := json.MarshalIndent(openApi, "", "    ")
	if err != nil {
		c.logger.WithError(err).Error("Failed to encode OpenApi document")
		return nil, err
	}
	jsonSchemaJSON = jsonBytes

	// Add a response:
	resFile := &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(openapiFileName),
		Content: proto.String(string(jsonSchemaJSON)),
	}

	return resFile, nil
}

func (c *Converter) createOpenApiDocument() (*OpenAPI, error) {
	openApi := OpenAPI{}
	if c.OpenApiFile != "" {
		fileContent, err := ioutil.ReadFile(c.OpenApiFile)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(fileContent, &openApi)
		if err != nil {
			return nil, err
		}
	}

	// OpenAPI version is mandatory
	if openApi.OpenAPI == "" {
		openApi.OpenAPI = "3.0.0"
	}

	// Paths is mandatory
	if openApi.Paths == nil {
		openApi.Paths = []byte("{}")
	}
	return &openApi, nil
}

func (c *Converter) calcGeneratorPlan(req *plugin.CodeGeneratorRequest) error {
	for _, file := range req.GetProtoFile() {
		packageName := file.GetPackage()
		for _, enum := range file.GetEnumType() {
			err := c.addToGeneratorPlan(c.outputFile, packageName, nil,
				nil, []*descriptor.EnumDescriptorProto{enum})
			if err != nil {
				return err
			}
		}
		for _, msg := range file.GetMessageType() {
			err := c.addToGeneratorPlan(c.outputFile, packageName, nil,
				[]*descriptor.DescriptorProto{msg}, nil)
			if err != nil {
				return err
			}
		}
	}

	c.logger.WithField("plan", c.generatorPlan).Debug("Generator Plan")
	return nil
}

func (c *Converter) addToGeneratorPlan(
	jsonFileName string,
	protoPackage string,
	parentTypeInfo *protoTypeInfo,
	childMessages []*descriptor.DescriptorProto,
	childEnums []*descriptor.EnumDescriptorProto) error {

	if len(childEnums) > 0 {
		for _, childEnum := range childEnums {
			typeInfo := NewProtoTypeInfoForEnum(jsonFileName, protoPackage, parentTypeInfo, childEnum)
			err := c.generatorPlan.Put(typeInfo)
			if err != nil {
				return err
			}
		}
	}

	if len(childMessages) > 0 {
		for _, childMsg := range childMessages {
			typeInfo := NewProtoTypeInfoForMsg(jsonFileName, protoPackage, parentTypeInfo, childMsg)
			err := c.generatorPlan.Put(typeInfo)
			if err != nil {
				return err
			}

			// recurse
			err = c.addToGeneratorPlan(
				jsonFileName, protoPackage, typeInfo, childMsg.GetNestedType(), childMsg.GetEnumType())
			if err != nil {
				return err
			}
		}
	}
	return nil
}
