package testdata

const ArrayOfObjects = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
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
                "$ref": "#/definitions/ArrayOfObjects.RepeatedPayload"
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
    "additionalProperties": true,
    "oneOf": [
        {
            "type": "null"
        },
        {
            "type": "object"
        }
    ],
    "definitions": {
        "ArrayOfObjects.RepeatedPayload": {
            "$schema": "http://json-schema.org/draft-04/schema#",
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
                    "$ref": "#/definitions/ArrayOfObjects.RepeatedPayload.Topology"
                }
            },
            "additionalProperties": true,
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
                },
                {
                    "type": "null"
                }
            ]
        }
    }
}`
