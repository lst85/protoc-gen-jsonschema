package converter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/alecthomas/jsonschema"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/sirupsen/logrus"
)

// Converter is everything you need to convert protos to JSONSchemas:
type Converter struct {
	// Allow NULL values for all properties. By default, JSONSchemas will reject NULL values.
	AllowNullValues bool
	// Allow additional properties. JSONSchemas will allow extra parameters, that are not specified in the schema.
	AllowAdditionalProperties bool
	// If the parameter is not set (default) the JSONSchema will allow both string and integers for 64 bit integers.
	// If it is set only integers are allowed.
	// The canonical JSON encoding of Proto3 converts int64, fixed64, uint6 to JSON strings.
	// When decoding JSON to ProtoBuf both numbers and strings are accepted.
	DisallowBigIntsAsStrings bool
	// Allow both enum names and integer values.
	AllowNumericEnumValues bool
	// Generate an OpenAPI v.3 file instead of JSONSchema file(s).
	// All ProtoBuf types (messages, enums, etc.) will be converted to their JSONSchema equivalent and added to the
	// components/schemas section of the OpenAPI document.
	// NOTE: The generator currently ignores gRPC service definitions. The paths section of the generated OpenAPI
	// document will be emtpy.
	GenerateOpenApi bool
	// Path to an OpenAPI file that will be merged with the generated schemas.
	// This parameter has only an effect when the parameter open_api is set.
	OpenApiFile string
	// Create a single file instead of multiple files.
	// When JSONSchema mode is enabled (parameter open_api is not set) a single JSONSchema will be generated with the
	// given filename.
	// When OpenAPI mode is enabled (parameter open_api is set) this parameter is implicitly enabled and the
	// default filename is "openapi.json".
	SingleOutputFile string
	// If the parameter is set the field names from the ProtoBuf definition are used in the JSONSchema.
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
				c.SingleOutputFile = paramSplit[1]
			}
		case parameter == "open_api":
			c.GenerateOpenApi = true
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
	if c.GenerateOpenApi && c.SingleOutputFile == "" {
		c.SingleOutputFile = "openapi.json"
	}
}

func (c *Converter) convert(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	c.sourceInfo = newSourceCodeInfo(req.GetProtoFile())
	res := &plugin.CodeGeneratorResponse{}

	err := c.calcGeneratorPlan(req)
	if err != nil {
		return nil, err
	}

	for _, fName := range c.generatorPlan.GetAllTargetFilenames() {
		converted, err := c.convertFile(fName)
		if err != nil {
			c.logger.WithField("file", fName).WithField("error", err).
				Error("Failed to create JSON schema file")
			res.Error = proto.String(fmt.Sprintf("Failed to create %s: %v", fName, err))
			return res, err
		}
		res.File = append(res.File, converted)
	}
	return res, nil
}

func (c *Converter) convertFile(jsonSchemaFileName string) (*plugin.CodeGeneratorResponse_File, error) {
	typeInfos := c.generatorPlan.GetAllForTargetFilename(jsonSchemaFileName)

	openApi, err := c.createOpenApiDocument()
	if err != nil {
		return nil, err
	}

	c.logger.WithField("file", jsonSchemaFileName).Debug("Creating JSON schema file ...")

	definitions := jsonschema.Definitions{}
	jsonSchema := &jsonschema.Schema{
		Definitions: definitions,
	}

	for _, typeInfo := range typeInfos {
		if typeInfo.GenerateAtTopLevel() {
			if jsonSchema.Type != nil {
				return nil, errors.New("Error while creating JSON Schema " + jsonSchemaFileName +
					". JSON schema can only have one root type.")
			}

			jsonType, err := c.convertProtoType(typeInfo)
			if err != nil {
				return nil, err
			}
			jsonSchema.Type = jsonType
		} else {
			jsonType, err := c.convertProtoType(typeInfo)
			if err != nil {
				return nil, err
			}

			typeId := typeInfo.GetProtoFQTypeName()
			if c.GenerateOpenApi {
				if openApi.Components.Schemas == nil {
					openApi.Components.Schemas = Schemas{}
				}
				openApi.Components.Schemas[typeId] = jsonType
			} else {
				jsonSchema.Definitions[typeId] = jsonType
			}
		}
	}

	// Marshal the result into a JSON-Schema or OpenAPI document:
	var jsonSchemaJSON []byte = nil
	if c.GenerateOpenApi {
		jsonBytes, err := json.MarshalIndent(openApi, "", "    ")
		if err != nil {
			c.logger.WithError(err).Error("Failed to encode OpenApi document")
			return nil, err
		}
		jsonSchemaJSON = jsonBytes
	} else {
		// JSON schema must have a top level type. When single file mode is enabled all types are generated in the
		// definitions section and thus we add an artifical top-level type to the schema.
		if jsonSchema.Type == nil {
			jsonSchema.Type = &jsonschema.Type{
				Type: "nil",
			}
		}
		jsonBytes, err := json.MarshalIndent(jsonSchema, "", "    ")
		if err != nil {
			c.logger.WithError(err).Error("Failed to encode jsonSchema")
			return nil, err
		}
		jsonSchemaJSON = jsonBytes
	}

	// Add a response:
	resFile := &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(jsonSchemaFileName),
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
			err := c.addToGeneratorPlan(c.getJsonFileName(packageName, enum.GetName()), packageName, nil,
				nil, []*descriptor.EnumDescriptorProto{enum})
			if err != nil {
				return err
			}
		}
		for _, msg := range file.GetMessageType() {
			err := c.addToGeneratorPlan(c.getJsonFileName(packageName, msg.GetName()), packageName, nil,
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

	jsonTopLevel := parentTypeInfo == nil && c.SingleOutputFile == ""

	if len(childEnums) > 0 {
		for _, childEnum := range childEnums {
			typeInfo := NewProtoTypeInfoForEnum(jsonFileName, jsonTopLevel, protoPackage, parentTypeInfo, childEnum)
			err := c.generatorPlan.Put(typeInfo)
			if err != nil {
				return err
			}
		}
	}

	if len(childMessages) > 0 {
		for _, childMsg := range childMessages {
			typeInfo := NewProtoTypeInfoForMsg(jsonFileName, jsonTopLevel, protoPackage, parentTypeInfo, childMsg)
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

func (c *Converter) getJsonFileName(packageName string, topLevelTypeName string) string {
	if c.SingleOutputFile != "" {
		return c.SingleOutputFile
	}

	res := packageName
	if packageName != "" {
		res += "."
	}
	return res + topLevelTypeName + ".json"
}
