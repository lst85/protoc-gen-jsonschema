package testdata

const AdditionalProperties = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "properties": {
        "key": {
            "type": "string"
        },
        "message": {
            "$ref": "#/definitions/AdditionalProperties.NestedMessage",
            "additionalProperties": {},
            "type": "object"
        },
        "repeatedPrimitive": {
            "items": {
                "type": "integer"
            },
            "type": "array"
        },
        "repeatedMessage": {
            "items": {
                "$ref": "#/definitions/AdditionalProperties.NestedMessage"
            },
            "type": "array"
        }
    },
    "additionalProperties": {},
    "type": "object",
    "definitions": {
        "AdditionalProperties.NestedMessage": {
            "$schema": "http://json-schema.org/draft-04/schema#",
            "properties": {
                "nestedKey": {
                    "type": "string"
                }
            },
            "additionalProperties": {},
            "type": "object"
        }
    }
}`
