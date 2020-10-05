package testdata

const ArrayOfEnums = `{
    "components": {
        "schemas": {
            "ArrayOfEnums": {
                "properties": {
                    "description": {
                        "type": "string"
                    },
                    "stuff": {
                        "items": {
                            "$ref": "#/components/schemas/ArrayOfEnums.inline"
                        },
                        "type": "array"
                    }
                },
                "type": "object"
            },
            "ArrayOfEnums.inline": {
                "enum": [
                    "FOO",
                    "BAR",
                    "FIZZ",
                    "BUZZ"
                ],
                "type": "string"
            }
        }
    },
    "openapi": "3.0.0",
    "paths": {}
}`
