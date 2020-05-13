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
            "additionalProperties": true,
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
    "additionalProperties": true,
    "type": "object",
    "definitions": {
        "Enumception.FailureModes": {
            "$schema": "http://json-schema.org/draft-04/schema#",
            "enum": [
                "RECURSION_ERROR",
                "SYNTAX_ERROR"
            ],
            "type": "string"
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
    "additionalProperties": true,
    "type": "object",
    "definitions": {
        "ImportedMessage.Topology": {
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

const EnumceptionImportedEnum = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "enum": [
        "VALUE_0",
        "VALUE_1",
        "VALUE_2",
        "VALUE_3"
    ],
    "type": "string"
}`
