package testdata

const AdditionalProperties = `{
    "components": {
        "schemas": {
            "AdditionalProperties": {
                "properties": {
                    "key": {
                        "type": "string"
                    },
                    "message": {
                        "$ref": "#/components/schemas/AdditionalProperties.NestedMessage",
                        "additionalProperties": {}
                    },
                    "repeatedPrimitive": {
                        "items": {
                            "type": "integer"
                        },
                        "type": "array"
                    },
                    "repeatedMessage": {
                        "items": {
                            "$ref": "#/components/schemas/AdditionalProperties.NestedMessage"
                        },
                        "type": "array"
                    }
                },
                "additionalProperties": {},
                "type": "object"
            },
            "AdditionalProperties.NestedMessage": {
                "properties": {
                    "nestedKey": {
                        "type": "string"
                    }
                },
                "additionalProperties": {},
                "type": "object"
            }
        }
    },
    "openapi": "3.0.0",
    "paths": {}
}`
