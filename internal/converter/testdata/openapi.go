package testdata

const OpenApi = `{
    "components": {
        "schemas": {
            "OpenApi": {
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
                        "$ref": "#/components/schemas/OpenApi.FailureModes"
                    },
                    "importedEnum": {
                        "$ref": "#/components/schemas/TopLevelEnum"
                    }
                },
                "type": "object"
            },
            "OpenApi.FailureModes": {
                "enum": [
                    "RECURSION_ERROR",
                    "SYNTAX_ERROR"
                ],
                "type": "string"
            },
            "TopLevelEnum": {
                "enum": [
                    "VALUE_0",
                    "VALUE_1",
                    "VALUE_2",
                    "VALUE_3"
                ],
                "type": "string"
            }
        }
    },
    "info": {
        "title": "My OpenAPI file",
        "version": "1.0"
    },
    "openapi": "3.0.0",
    "paths": {},
    "servers": [
        {
            "url": "/services"
        }
    ]
}`
