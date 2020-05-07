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
	SingleFile                   bool
	UseProtoAndJSONFieldnames    bool
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
		switch parameter {
		case "": // ignore empty parameter
		case "allow_null_values":
			c.AllowNullValues = true
		case "debug":
			c.logger.SetLevel(logrus.DebugLevel)
		case "disallow_additional_properties":
			c.DisallowAdditionalProperties = true
		case "disallow_bigints_as_strings":
			c.DisallowBigIntsAsStrings = true
		case "disallow_numeric_enum_values":
			c.DisallowNumericEnumValues = true
		case "single_file":
			c.SingleFile = true
		case "open_api_conform":
			c.OpenApiConform = true
		case "proto_and_json_fieldnames":
			c.UseProtoAndJSONFieldnames = true
		default:
			c.logger.WithField("parameter", parameter).Warn("Unknown parameter")
		}
	}
}

func (c *Converter) convertFile(jsonSchemaFileName string) (*plugin.CodeGeneratorResponse_File, error) {
	typeInfos := c.generatorPlan.GetForJsonFilename(jsonSchemaFileName)

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

			jsonType, err := c.recursiveConvertMessageType(typeInfo)
			if err != nil {
				return nil, err
			}
			jsonSchema.Type = jsonType
		} else {
			jsonType, err := c.recursiveConvertMessageType(typeInfo)
			if err != nil {
				return nil, err
			}

			// need to give that schema an ID
			typeId := typeInfo.GetProtoFQNName()
			if jsonType.Extras == nil {
				jsonType.Extras = make(map[string]interface{})
			}
			jsonType.Extras["id"] = typeId
			jsonSchema.Definitions[typeId] = jsonType
		}
	}

	// Marshal the JSON-Schema into JSON:
	if jsonSchema.Type == nil {
		jsonSchema.Type = &jsonschema.Type{}
	}
	jsonSchemaJSON, err := json.MarshalIndent(jsonSchema, "", "    ")
	if err != nil {
		c.logger.WithError(err).Error("Failed to encode jsonSchema")
		return nil, err
	}

	// Add a response:
	resFile := &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(jsonSchemaFileName),
		Content: proto.String(string(jsonSchemaJSON)),
	}

	return resFile, nil
}

func (c *Converter) convert(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	c.sourceInfo = newSourceCodeInfo(req.GetProtoFile())
	res := &plugin.CodeGeneratorResponse{}

	err := c.calcGeneratorPlan(req)
	if err != nil {
		return nil, err
	}

	for _, fName := range c.generatorPlan.GetAllJsonFilenames() {
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

	jsonTopLevel := parentTypeInfo == nil && !c.SingleFile

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
	if c.SingleFile {
		return "schema.json"
	}

	res := packageName
	if packageName != "" {
		res += "."
	}
	return res + topLevelTypeName + ".json"
}
