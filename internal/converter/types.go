package converter

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/alecthomas/jsonschema"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/iancoleman/orderedmap"
	"github.com/xeipuuv/gojsonschema"
)

func (c *Converter) convertProtoType(typeInfo *protoTypeInfo) (*jsonschema.Type, error) {
	if typeInfo.IsProtoMessage() {
		return c.convertMessageType(typeInfo, typeInfo.GetProtoMessage())
	} else if typeInfo.IsProtoEnum() {
		return c.convertEnumType(typeInfo.GetProtoEnum())
	} else {
		return nil, errors.New("unknown ProtoBuf type, expecting either message or enum")
	}
}

// Converts a proto "ENUM" into a JSON-Schema:
func (c *Converter) convertEnumType(enum *descriptor.EnumDescriptorProto) (*jsonschema.Type, error) {

	// Prepare a new jsonschema.Type for our eventual return value:
	jsonSchemaType := c.createNewJsonType()

	// Generate a description from src comments (if available)
	if src := c.sourceInfo.GetEnum(enum); src != nil {
		jsonSchemaType.Description = formatDescription(src)
	}

	c.setJsonTypeForEnum(jsonSchemaType)

	// Add the allowed values:
	for _, enumValue := range enum.Value {
		jsonSchemaType.Enum = append(jsonSchemaType.Enum, enumValue.Name)
		if c.AllowNumericEnumValues {
			jsonSchemaType.Enum = append(jsonSchemaType.Enum, enumValue.Number)
		}
	}

	return jsonSchemaType, nil
}

// Converts a proto "message" into a JSON-Schema:
func (c *Converter) convertMessageType(typeInfo *protoTypeInfo, msg *descriptor.DescriptorProto) (*jsonschema.Type, error) {
	// Prepare a new jsonschema:
	jsonSchemaType := c.createNewJsonType()
	jsonSchemaType.Properties = orderedmap.New()

	// Generate a description from src comments (if available)
	if src := c.sourceInfo.GetMessage(msg); src != nil {
		jsonSchemaType.Description = formatDescription(src)
	}

	// Optionally allow NULL values:
	if c.AllowNullValues {
		jsonSchemaType.OneOf = []*jsonschema.Type{
			{Type: gojsonschema.TYPE_NULL},
			{Type: gojsonschema.TYPE_OBJECT},
		}
	} else {
		jsonSchemaType.Type = gojsonschema.TYPE_OBJECT
	}

	jsonSchemaType.AdditionalProperties = c.getAdditionalPropertiesValue()

	c.logger.WithField("message_str", proto.MarshalTextString(msg)).Trace("Converting message")
	for _, fieldDesc := range msg.GetField() {
		recursedJSONSchemaType, err := c.convertField(typeInfo, fieldDesc)
		if err != nil {
			c.logger.WithError(err).WithField("field_name", fieldDesc.GetName()).WithField("message_name", msg.GetName()).Error("Failed to convert field")
			return nil, err
		}
		c.logger.WithField("field_name", fieldDesc.GetName()).WithField("type", recursedJSONSchemaType.Type).Trace("Converted field")

		if c.UseProtoFieldNames {
			jsonSchemaType.Properties.Set(fieldDesc.GetName(), recursedJSONSchemaType)
		} else {
			jsonSchemaType.Properties.Set(fieldDesc.GetJsonName(), recursedJSONSchemaType)
		}
	}

	if len(jsonSchemaType.Properties.Keys()) == 0 {
		// remove empty properties to clean the final output as clean as possible
		jsonSchemaType.Properties = nil
	}

	return jsonSchemaType, nil
}

// Convert a proto "field" (essentially a type-switch with some recursion):
func (c *Converter) convertField(typeInfo *protoTypeInfo, desc *descriptor.FieldDescriptorProto) (*jsonschema.Type, error) {

	// Prepare a new jsonschema.Type for our eventual return value:
	jsonSchemaType := &jsonschema.Type{}

	// Generate a description from src comments (if available)
	if src := c.sourceInfo.GetField(desc); src != nil {
		jsonSchemaType.Description = formatDescription(src)
	}

	// Switch the types, and pick a JSONSchema equivalent:
	switch desc.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT:
		if c.AllowNullValues {
			jsonSchemaType.OneOf = []*jsonschema.Type{
				{Type: gojsonschema.TYPE_NULL},
				{Type: gojsonschema.TYPE_NUMBER},
			}
		} else {
			jsonSchemaType.Type = gojsonschema.TYPE_NUMBER
		}

	case descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_SINT32:
		if c.AllowNullValues {
			jsonSchemaType.OneOf = []*jsonschema.Type{
				{Type: gojsonschema.TYPE_NULL},
				{Type: gojsonschema.TYPE_INTEGER},
			}
		} else {
			jsonSchemaType.Type = gojsonschema.TYPE_INTEGER
		}

	case descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_SINT64:
		jsonSchemaType.OneOf = append(jsonSchemaType.OneOf, &jsonschema.Type{Type: gojsonschema.TYPE_INTEGER})
		if !c.DisallowBigIntsAsStrings {
			jsonSchemaType.OneOf = append(jsonSchemaType.OneOf, &jsonschema.Type{Type: gojsonschema.TYPE_STRING})
		}
		if c.AllowNullValues {
			jsonSchemaType.OneOf = append(jsonSchemaType.OneOf, &jsonschema.Type{Type: gojsonschema.TYPE_NULL})
		}

	case descriptor.FieldDescriptorProto_TYPE_STRING,
		descriptor.FieldDescriptorProto_TYPE_BYTES:
		if c.AllowNullValues {
			jsonSchemaType.OneOf = []*jsonschema.Type{
				{Type: gojsonschema.TYPE_NULL},
				{Type: gojsonschema.TYPE_STRING},
			}
		} else {
			jsonSchemaType.Type = gojsonschema.TYPE_STRING
		}

	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		fieldType := c.generatorPlan.LookupType(desc.GetTypeName())
		if fieldType == nil {
			return nil, fmt.Errorf("no such message type: %s", desc.GetTypeName())
		}
		jsonSchemaType.Ref = c.getJsonRefValue(typeInfo, fieldType)

	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		if c.AllowNullValues {
			jsonSchemaType.OneOf = []*jsonschema.Type{
				{Type: gojsonschema.TYPE_NULL},
				{Type: gojsonschema.TYPE_BOOLEAN},
			}
		} else {
			jsonSchemaType.Type = gojsonschema.TYPE_BOOLEAN
		}

	case descriptor.FieldDescriptorProto_TYPE_GROUP,
		descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		jsonSchemaType.Type = gojsonschema.TYPE_OBJECT

	default:
		return nil, fmt.Errorf("unrecognized field type: %s", desc.GetType().String())
	}

	// Recurse array of primitive types:
	if desc.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED && jsonSchemaType.Type != gojsonschema.TYPE_OBJECT {
		jsonSchemaType.Items = &jsonschema.Type{}

		// Copy all relevant properties to the items section:
		jsonSchemaType.Items.Enum = jsonSchemaType.Enum
		jsonSchemaType.Enum = nil

		jsonSchemaType.Items.Type = jsonSchemaType.Type
		jsonSchemaType.Type = ""

		jsonSchemaType.Items.OneOf = jsonSchemaType.OneOf
		jsonSchemaType.OneOf = nil

		jsonSchemaType.Items.Ref = jsonSchemaType.Ref
		jsonSchemaType.Ref = ""

		// Set the type to array:
		if c.AllowNullValues {
			jsonSchemaType.OneOf = []*jsonschema.Type{
				{Type: gojsonschema.TYPE_NULL},
				{Type: gojsonschema.TYPE_ARRAY},
			}
		} else {
			jsonSchemaType.Type = gojsonschema.TYPE_ARRAY
			jsonSchemaType.OneOf = []*jsonschema.Type{}
		}
		return jsonSchemaType, nil
	}

	// Recurse nested objects / arrays of objects (if necessary):
	if jsonSchemaType.Type == gojsonschema.TYPE_OBJECT {

		fieldType := c.generatorPlan.LookupType(desc.GetTypeName())
		if fieldType == nil {
			return nil, fmt.Errorf("no such message type: %s", desc.GetTypeName())
		}

		if !fieldType.IsProtoMessage() {
			return nil, fmt.Errorf("expecting type %s to be a ProtoBuf message", fieldType.GetProtoFQTypeName())
		}
		fieldMsg := fieldType.GetProtoMessage()

		// Maps, arrays, and objects are structured in different ways:
		switch {

		case fieldMsg.Options.GetMapEntry():
			c.logger.
				WithField("field_name", fieldMsg.GetName()).
				WithField("msg_name", typeInfo.GetProtoFQTypeName()).
				Tracef("Is a map")

			// Recurse the recordType:
			recursedJSONSchemaType, err := c.convertProtoType(fieldType)
			if err != nil {
				return nil, err
			}

			// Make sure we have a "value":
			value, valuePresent := recursedJSONSchemaType.Properties.Get("value")
			if !valuePresent {
				return nil, fmt.Errorf("unable to find 'value' property of MAP type")
			}

			// Marshal the "value" properties to JSON (because that's how we can pass on AdditionalProperties):
			additionalPropertiesJSON, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}
			jsonSchemaType.AdditionalProperties = additionalPropertiesJSON

		// Arrays:
		case desc.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED:
			jsonSchemaType.Type = gojsonschema.TYPE_ARRAY
			jsonSchemaType.Items = &jsonschema.Type{}
			jsonSchemaType.Items.Ref = c.getJsonRefValue(typeInfo, fieldType)

		// Objects:
		default:
			jsonSchemaType.Ref = c.getJsonRefValue(typeInfo, fieldType)
			jsonSchemaType.Type = "" // If $ref is set the property type must not be set
			jsonSchemaType.AdditionalProperties = c.getAdditionalPropertiesValue()
		}

		// Optionally allow NULL values:
		if c.AllowNullValues {
			jsonSchemaType.OneOf = []*jsonschema.Type{
				{Type: gojsonschema.TYPE_NULL},
				{Type: jsonSchemaType.Type},
			}
			jsonSchemaType.Type = ""
		}
	}

	return jsonSchemaType, nil
}

func formatDescription(sl *descriptor.SourceCodeInfo_Location) string {
	var lines []string
	for _, str := range sl.GetLeadingDetachedComments() {
		if s := strings.TrimSpace(str); s != "" {
			lines = append(lines, s)
		}
	}
	if s := strings.TrimSpace(sl.GetLeadingComments()); s != "" {
		lines = append(lines, s)
	}
	if s := strings.TrimSpace(sl.GetTrailingComments()); s != "" {
		lines = append(lines, s)
	}
	return strings.Join(lines, "\n\n")
}

// disallowAdditionalProperties will prevent validation where extra fields are found (outside of the schema).
// The specification of additionalProperties is different in OpenAPI compared to JSON schema:
// 1. Default when additionalProperties is not provided:
//		JSON schema: all additional properties allowed
//		OpenAPI: no additional properties allowed
// 2. Some OpenAPI implementations do not allow boolean values.
func (c *Converter) getAdditionalPropertiesValue() []byte {
	if c.GenerateOpenApi {
		if c.AllowAdditionalProperties {
			return []byte("{}")
		}
		return nil
	}

	if c.AllowAdditionalProperties {
		return []byte("true")
	}
	return []byte("false")
}

func (c *Converter) setJsonTypeForEnum(jsonSchemaType *jsonschema.Type) {
	types := []string{gojsonschema.TYPE_STRING}
	if c.AllowNumericEnumValues {
		types = append(types, gojsonschema.TYPE_INTEGER)
	}
	if c.AllowNullValues {
		types = append(types, gojsonschema.TYPE_NULL)
	}

	if len(types) == 1 {
		jsonSchemaType.Type = types[0]
	} else {
		for i := range types {
			jsonSchemaType.OneOf = append(jsonSchemaType.OneOf, &jsonschema.Type{Type: types[i]})
		}
	}
}

func (c *Converter) getJsonRefValue(contextType *protoTypeInfo, targetType *protoTypeInfo) string {
	ref := ""
	if contextType.GetTargetFileName() != targetType.GetTargetFileName() {
		ref += targetType.GetTargetFileName()
	}
	ref += "#"
	if !targetType.GenerateAtTopLevel() {
		if c.GenerateOpenApi {
			ref += "/components/schemas/"
		} else {
			ref += "/definitions/"
		}
		ref += targetType.GetProtoFQTypeName()
	}
	return ref
}

func (c *Converter) createNewJsonType() *jsonschema.Type {
	t := jsonschema.Type{}
	if !c.GenerateOpenApi {
		t.Version = jsonschema.Version
	}
	return &t
}
