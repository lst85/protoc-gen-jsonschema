package testdata

const Maps = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "properties": {
        "mapOfStrings": {
            "additionalProperties": {
                "type": "string"
            },
            "type": "object"
        },
        "mapOfInts": {
            "additionalProperties": {
                "type": "integer"
            },
            "type": "object"
        },
        "mapOfMessages": {
            "additionalProperties": {
                "$ref": "#/definitions/Maps.PayloadMessage",
                "additionalProperties": true,
                "type": "object"
            },
            "type": "object"
        }
    },
    "additionalProperties": true,
    "type": "object",
    "definitions": {
        "Maps.MapOfIntsEntry": {
            "$schema": "http://json-schema.org/draft-04/schema#",
            "properties": {
                "key": {
                    "type": "string"
                },
                "value": {
                    "type": "integer"
                }
            },
            "additionalProperties": true,
            "type": "object"
        },
        "Maps.MapOfMessagesEntry": {
            "$schema": "http://json-schema.org/draft-04/schema#",
            "properties": {
                "key": {
                    "type": "string"
                },
                "value": {
                    "$ref": "#/definitions/Maps.PayloadMessage",
                    "additionalProperties": true,
                    "type": "object"
                }
            },
            "additionalProperties": true,
            "type": "object"
        },
        "Maps.MapOfStringsEntry": {
            "$schema": "http://json-schema.org/draft-04/schema#",
            "properties": {
                "key": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            },
            "additionalProperties": true,
            "type": "object"
        },
        "Maps.PayloadMessage": {
            "$schema": "http://json-schema.org/draft-04/schema#",
            "properties": {
                "name": {
                    "type": "string"
                },
                "topology": {
                    "$ref": "#/definitions/Maps.PayloadMessage.Topology"
                }
            },
            "additionalProperties": true,
            "type": "object"
        },
        "Maps.PayloadMessage.Topology": {
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
