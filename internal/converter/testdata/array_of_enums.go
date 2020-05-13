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
