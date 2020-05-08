package testdata

const ArrayOfEnums = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "properties": {
        "description": {
            "type": "string"
        },
        "stuff": {
            "$ref": "#/definitions/ArrayOfEnums.inline",
            "items": {},
            "type": "array"
        }
    },
    "additionalProperties": true,
    "type": "object",
    "definitions": {
        "ArrayOfEnums.inline": {
            "$schema": "http://json-schema.org/draft-04/schema#",
            "enum": [
                "FOO",
                0,
                "BAR",
                1,
                "FIZZ",
                2,
                "BUZZ",
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
        }
    }
}`
