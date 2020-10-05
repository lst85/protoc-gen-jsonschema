package testdata

const Enumception = `{
    "components": {
        "schemas": {
            "Enumception": {
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
                        "$ref": "#/components/schemas/Enumception.FailureModes"
                    },
                    "payload": {
                        "$ref": "#/components/schemas/ImportedMessage",
                        "additionalProperties": {}
                    },
                    "payloads": {
                        "items": {
                            "$ref": "#/components/schemas/ImportedMessage"
                        },
                        "type": "array"
                    },
                    "importedEnum": {
                        "$ref": "#/components/schemas/ImportedEnum"
                    },
                    "importedEnums": {
                        "items": {
                            "$ref": "#/components/schemas/ImportedEnum"
                        },
                        "type": "array"
                    }
                },
                "additionalProperties": {},
                "type": "object"
            },
            "Enumception.FailureModes": {
                "enum": [
                    "RECURSION_ERROR",
                    "SYNTAX_ERROR"
                ],
                "type": "string"
            },
            "ImportedEnum": {
                "enum": [
                    "VALUE_0",
                    "VALUE_1",
                    "VALUE_2",
                    "VALUE_3"
                ],
                "type": "string"
            },
            "ImportedMessage": {
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
                        "$ref": "#/components/schemas/ImportedMessage.Topology"
                    }
                },
                "additionalProperties": {},
                "type": "object"
            },
            "ImportedMessage.Topology": {
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
