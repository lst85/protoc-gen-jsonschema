package testdata

const SelfReference = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "properties": {
        "name": {
            "type": "string"
        },
        "bar": {
            "items": {
                "$ref": "#"
            },
            "type": "array"
        }
    },
    "additionalProperties": true,
    "type": "object"
}`
