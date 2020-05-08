package testdata

const Enumception = `{
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
        "failureMode": {
            "$ref": "#/definitions/Enumception.FailureModes"
        },
        "payload": {
            "$ref": "samples.enumception.ImportedMessage.json#",
            "additionalProperties": false,
            "type": "object"
        },
        "payloads": {
            "items": {
                "$ref": "samples.enumception.ImportedMessage.json#"
            },
            "type": "array"
        },
        "importedEnum": {
            "$ref": "samples.enumception.ImportedEnum.json#"
        }
    },
    "additionalProperties": false,
    "type": "object",
    "definitions": {
        "Enumception.FailureModes": {
            "$schema": "http://json-schema.org/draft-04/schema#",
            "enum": [
                "RECURSION_ERROR",
                0,
                "SYNTAX_ERROR",
                1
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

const EnumceptionImportedMessage = `{
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
            "$ref": "#/definitions/ImportedMessage.Topology"
        }
    },
    "additionalProperties": false,
    "type": "object",
    "definitions": {
        "ImportedMessage.Topology": {
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

const EnumceptionImportedEnum = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "enum": [
        "VALUE_0",
        0,
        "VALUE_1",
        1,
        "VALUE_2",
        2,
        "VALUE_3",
        3
    ],
    "oneOf": [
        {
            "type": "string"
        },
        {
            "type": "integer"
        }
    ]
}`
