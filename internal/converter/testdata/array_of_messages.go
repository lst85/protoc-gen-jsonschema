package testdata

const ArrayOfMessages = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "properties": {
        "description": {
            "type": "string"
        },
        "payload": {
            "items": {
                "$ref": "#/definitions/ArrayOfMessages.PayloadMessage"
            },
            "type": "array"
        }
    },
    "additionalProperties": true,
    "type": "object",
    "definitions": {
        "ArrayOfMessages.PayloadMessage": {
            "$schema": "http://json-schema.org/draft-04/schema#",
            "properties": {
                "name": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "rating": {
                    "type": "number"
                },
                "complete": {
                    "type": "boolean"
                },
                "topology": {
                    "$ref": "#/definitions/ArrayOfMessages.PayloadMessage.Topology"
                }
            },
            "additionalProperties": true,
            "type": "object"
        },
        "ArrayOfMessages.PayloadMessage.Topology": {
            "$schema": "http://json-schema.org/draft-04/schema#",
            "enum": [
                "FLAT",
                0,
                "NESTED_OBJECT",
                1,
                "NESTED_MESSAGE",
                2,
                "ARRAY_OF_TYPE",
                3,
                "ARRAY_OF_OBJECT",
                4,
                "ARRAY_OF_MESSAGE",
                5
            ],
            "oneOf": [
                {
                    "type": "string"
                },
                {
                    "type": "integer"
                }
            ]
        }
    }
}`
