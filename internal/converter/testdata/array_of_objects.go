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
    "additionalProperties": false,
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
            "additionalProperties": false,
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
}`
