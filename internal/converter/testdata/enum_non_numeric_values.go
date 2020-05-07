package testdata

const EnumNoNumericValuesMsg = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "properties": {
        "description": {
            "type": "string"
        },
        "stuff": {
            "enum": [
                "FOO",
                "BAR",
                "FIZZ",
                "BUZZ"
            ],
            "type": "string"
        },
        "nestedEnum": {
            "type": "string"
        },
        "topLevelEnum": {
            "type": "string"
        }
    },
    "additionalProperties": true,
    "type": "object"
}`

const EnumNoNumericValuesMsg2 = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "additionalProperties": true,
    "type": "object"
}`

const EnumNoNumericValuesTopLevelEnum = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "enum": [
        "FRR",
        "FRA"
    ],
    "type": "string"
}`
