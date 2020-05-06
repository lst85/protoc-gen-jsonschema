package testdata

const AdditionalProperties = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "properties": {
        "key": {
            "type": "string"
        },
        "message": {
            "$ref": "#/definitions/samples.AdditionalProperties.NestedMessage",
            "additionalProperties": {},
            "type": "object"
        },
        "repeated_primitive": {
            "items": {
                "type": "integer"
            },
            "type": "array"
        },
        "repeated_message": {
            "items": {
                "$schema": "http://json-schema.org/draft-04/schema#",
                "$ref": "#/definitions/samples.AdditionalProperties.NestedMessage"
            },
            "type": "array"
        }
    },
    "additionalProperties": {},
    "type": "object",
    "definitions": {
        "samples.AdditionalProperties.NestedMessage": {
            "$schema": "http://json-schema.org/draft-04/schema#",
            "properties": {
                "nested_key": {
                    "type": "string"
                }
            },
            "additionalProperties": {},
            "type": "object",
            "id": "samples.AdditionalProperties.NestedMessage"
        }
    }
}`
