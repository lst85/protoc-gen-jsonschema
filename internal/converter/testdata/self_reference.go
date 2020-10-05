package testdata

const SelfReference = `{
    "components": {
        "schemas": {
            "SelfReference": {
                "properties": {
                    "name": {
                        "type": "string"
                    },
                    "bar": {
                        "items": {
                            "$ref": "#/components/schemas/SelfReference"
                        },
                        "type": "array"
                    }
                },
                "type": "object"
            }
        }
    },
    "openapi": "3.0.0",
    "paths": {}
}`
