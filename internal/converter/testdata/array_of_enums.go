package testdata

const ArrayOfEnums = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "properties": {
        "description": {
            "type": "string"
        },
        "stuff": {
            "items": {
                "$ref": "#/definitions/ArrayOfEnums.inline"
            },
            "type": "array"
        }
    },
    "additionalProperties": false,
    "type": "object",
    "definitions": {
        "ArrayOfEnums.inline": {
            "$schema": "http://json-schema.org/draft-04/schema#",
            "enum": [
                "FOO",
                "BAR",
                "FIZZ",
                "BUZZ"
            ],
            "type": "string"
        }
    }
}`
