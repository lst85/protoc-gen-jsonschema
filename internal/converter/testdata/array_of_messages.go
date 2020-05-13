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
    "additionalProperties": false,
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
            "additionalProperties": false,
            "type": "object"
        },
        "ArrayOfMessages.PayloadMessage.Topology": {
            "$schema": "http://json-schema.org/draft-04/schema#",
            "enum": [
                "FLAT",
                "NESTED_OBJECT",
                "NESTED_MESSAGE",
                "ARRAY_OF_TYPE",
                "ARRAY_OF_OBJECT",
                "ARRAY_OF_MESSAGE"
            ],
            "type": "string"
        }
    }
}`
