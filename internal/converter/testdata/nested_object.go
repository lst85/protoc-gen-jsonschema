package testdata

const NestedObject = `{
    "components": {
        "schemas": {
            "NestedObject": {
                "properties": {
                    "payload": {
                        "$ref": "#/components/schemas/NestedObject.NestedPayload"
                    },
                    "description": {
                        "type": "string"
                    }
                },
                "type": "object"
            },
            "NestedObject.NestedPayload": {
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
                        "$ref": "#/components/schemas/NestedObject.NestedPayload.Topology"
                    }
                },
                "type": "object"
            },
            "NestedObject.NestedPayload.Topology": {
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
