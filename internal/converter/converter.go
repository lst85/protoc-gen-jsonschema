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
	AllowNullValues              bool
	DisallowAdditionalProperties bool
	DisallowBigIntsAsStrings     bool
	DisallowNumericEnumValues    bool
	OpenApiConform               bool
	OpenApiFile                  string
	SingleOutputFile             string
	UseProtoFieldnames           bool
	logger                       *logrus.Logger
	sourceInfo                   *sourceCodeInfo
	generatorPlan                *generatorPlan
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
		case parameter == "debug":
			c.logger.SetLevel(logrus.DebugLevel)
		case parameter == "disallow_additional_properties":
			c.DisallowAdditionalProperties = true
		case parameter == "disallow_bigints_as_strings":
			c.DisallowBigIntsAsStrings = true
		case parameter == "disallow_numeric_enum_values":
			c.DisallowNumericEnumValues = true
		case strings.HasPrefix(parameter, "out_file"):
			paramSplit := strings.Split(parameter, "=")
			if len(paramSplit) != 2 || len(paramSplit[1]) == 0 {
				c.logger.WithField("parameter", parameter).
					Warn("Invalid parameter. Expecting out_file=<filename>")
			} else {
				c.SingleOutputFile = paramSplit[1]
			}
		case parameter == "open_api_conform":
			c.OpenApiConform = true
		case strings.HasPrefix(parameter, "open_api"):
			paramSplit := strings.Split(parameter, "=")
			if len(paramSplit) != 2 || len(paramSplit[1]) == 0 {
				c.logger.WithField("parameter", parameter).
					Warn("Invalid parameter. Expecting open_api=<filename>")
			} else {
				c.OpenApiFile = paramSplit[1]
			}
		case parameter == "proto_fieldnames":
			c.UseProtoFieldnames = true
		default:
			c.logger.WithField("parameter", parameter).Warn("Unknown parameter")
		}
	}

	if c.OpenApiFile != "" && c.SingleOutputFile == "" {
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

	c.logger.WithField("file", jsonSchemaFileName).Info("Creating JSON schema file ...")

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

			// need to give that schema an ID
			typeId := typeInfo.GetProtoFQTypeName()
			if c.OpenApiFile != "" {
				if openApi.Components.Schemas == nil {
					openApi.Components.Schemas = Schemas{}
				}
				openApi.Components.Schemas[typeId] = jsonType
			} else {
				jsonSchema.Definitions[typeId] = jsonType
			}
		}
	}

	var jsonSchemaJSON []byte = nil
	if c.OpenApiFile != "" {
		if openApi.OpenAPI == "" {
			openApi.OpenAPI = "3.0.0"
		}
		jsonBytes, err := json.MarshalIndent(openApi, "", "    ")
		if err != nil {
			c.logger.WithError(err).Error("Failed to encode OpenApi document")
			return nil, err
		}
		jsonSchemaJSON = jsonBytes
	} else {
		// Marshal the JSON-Schema into JSON:
		if jsonSchema.Type == nil {
			jsonSchema.Type = &jsonschema.Type{}
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

	c.logger.WithField("plan", c.generatorPlan).Info("Generator Plan")
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
