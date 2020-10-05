package testdata

const ArrayOfObjects = `{
    "components": {
        "schemas": {
            "ArrayOfObjects": {
                "properties": {
                    "description": {
                        "oneOf": [
                            {
                                "type": "null"
                            },
                            {
                                "type": "string"
                            }
                        ]
                    },
                    "payload": {
                        "items": {
                            "$ref": "#/components/schemas/ArrayOfObjects.RepeatedPayload"
                        },
                        "oneOf": [
                            {
                                "type": "null"
                            },
                            {
                                "type": "array"
                            }
                        ]
                    }
                },
                "oneOf": [
                    {
                        "type": "null"
                    },
                    {
                        "type": "object"
                    }
                ]
            },
            "ArrayOfObjects.RepeatedPayload": {
                "properties": {
                    "name": {
                        "oneOf": [
                            {
                                "type": "null"
                            },
                            {
                                "type": "string"
                            }
                        ]
                    },
                    "timestamp": {
                        "oneOf": [
                            {
                                "type": "null"
                            },
                            {
                                "type": "string"
                            }
                        ]
                    },
                    "id": {
                        "oneOf": [
                            {
                                "type": "null"
                            },
                            {
                                "type": "integer"
                            }
                        ]
                    },
                    "rating": {
                        "oneOf": [
                            {
                                "type": "null"
                            },
                            {
                                "type": "number"
                            }
                        ]
                    },
                    "complete": {
                        "oneOf": [
                            {
                                "type": "null"
                            },
                            {
                                "type": "boolean"
                            }
                        ]
                    },
                    "topology": {
                        "$ref": "#/components/schemas/ArrayOfObjects.RepeatedPayload.Topology"
                    }
                },
                "oneOf": [
                    {
                        "type": "null"
                    },
                    {
                        "type": "object"
                    }
                ]
            },
            "ArrayOfObjects.RepeatedPayload.Topology": {
                "enum": [
                    "FLAT",
                    "NESTED_OBJECT",
                    "NESTED_MESSAGE",
                    "ARRAY_OF_TYPE",
                    "ARRAY_OF_OBJECT",
                    "ARRAY_OF_MESSAGE"
                ],
                "oneOf": [
                    {
                        "type": "string"
                    },
                    {
                        "type": "null"
                    }
                ]
            }
        }
    },
    "openapi": "3.0.0",
    "paths": {}
}`
