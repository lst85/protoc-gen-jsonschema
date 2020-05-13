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
                "additionalProperties": false,
                "type": "object"
            },
            "type": "object"
        }
    },
    "additionalProperties": false,
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
            "additionalProperties": false,
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
                    "additionalProperties": false,
                    "type": "object"
                }
            },
            "additionalProperties": false,
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
            "additionalProperties": false,
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
            "additionalProperties": false,
            "type": "object"
        },
        "Maps.PayloadMessage.Topology": {
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
