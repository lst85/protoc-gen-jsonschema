package testdata

const NestedObject = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "properties": {
        "payload": {
            "$ref": "#/definitions/NestedObject.NestedPayload",
            "additionalProperties": false,
            "type": "object"
        },
        "description": {
            "type": "string"
        }
    },
    "additionalProperties": false,
    "type": "object",
    "definitions": {
        "NestedObject.NestedPayload": {
            "$schema": "http://json-schema.org/draft-04/schema#",
            "properties": {
                "name": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "string"
                },
                "id": {
                    "oneOf": [
                        {
                            "type": "integer"
                        }
                    ]
                },
                "largeValue": {
                    "oneOf": [
                        {
                            "type": "integer"
                        }
                    ]
                },
                "rating": {
                    "type": "number"
                },
                "complete": {
                    "type": "boolean"
                },
                "topology": {
                    "$ref": "#/definitions/NestedObject.NestedPayload.Topology"
                }
            },
            "additionalProperties": false,
            "type": "object"
        },
        "NestedObject.NestedPayload.Topology": {
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
