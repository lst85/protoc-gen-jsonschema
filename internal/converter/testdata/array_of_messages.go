package testdata

const ArrayOfMessages = `{
    "components": {
        "schemas": {
            "ArrayOfMessages": {
                "properties": {
                    "description": {
                        "type": "string"
                    },
                    "payload": {
                        "items": {
                            "$ref": "#/components/schemas/ArrayOfMessages.PayloadMessage"
                        },
                        "type": "array"
                    }
                },
                "type": "object"
            },
            "ArrayOfMessages.PayloadMessage": {
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
                        "$ref": "#/components/schemas/ArrayOfMessages.PayloadMessage.Topology"
                    }
                },
                "type": "object"
            },
            "ArrayOfMessages.PayloadMessage.Topology": {
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
    },
    "openapi": "3.0.0",
    "paths": {}
}`
