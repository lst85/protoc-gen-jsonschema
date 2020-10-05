package testdata

const Maps = `{
    "components": {
        "schemas": {
            "Maps": {
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
                            "$ref": "#/components/schemas/Maps.PayloadMessage"
                        },
                        "type": "object"
                    }
                },
                "type": "object"
            },
            "Maps.MapOfIntsEntry": {
                "properties": {
                    "key": {
                        "type": "string"
                    },
                    "value": {
                        "type": "integer"
                    }
                },
                "type": "object"
            },
            "Maps.MapOfMessagesEntry": {
                "properties": {
                    "key": {
                        "type": "string"
                    },
                    "value": {
                        "$ref": "#/components/schemas/Maps.PayloadMessage"
                    }
                },
                "type": "object"
            },
            "Maps.MapOfStringsEntry": {
                "properties": {
                    "key": {
                        "type": "string"
                    },
                    "value": {
                        "type": "string"
                    }
                },
                "type": "object"
            },
            "Maps.PayloadMessage": {
                "properties": {
                    "name": {
                        "type": "string"
                    },
                    "topology": {
                        "$ref": "#/components/schemas/Maps.PayloadMessage.Topology"
                    }
                },
                "type": "object"
            },
            "Maps.PayloadMessage.Topology": {
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
