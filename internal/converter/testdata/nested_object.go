package testdata

const NestedObject = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "properties": {
        "payload": {
            "$ref": "#/definitions/NestedObject.NestedPayload",
            "additionalProperties": true,
            "type": "object"
        },
        "description": {
            "type": "string"
        }
    },
    "additionalProperties": true,
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
                    "type": "integer"
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
            "additionalProperties": true,
            "type": "object"
        },
        "NestedObject.NestedPayload.Topology": {
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
