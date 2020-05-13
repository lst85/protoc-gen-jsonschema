package testdata

const EnumNumericValuesMsg = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "properties": {
        "description": {
            "type": "string"
        },
        "stuff": {
            "$ref": "#/definitions/EnumNumericValuesMsg.InlineEnum"
        },
        "nestedEnum": {
            "$ref": "samples.subpackage.EnumNumericValuesMsg2.json#/definitions/EnumNumericValuesMsg2.NestedEnum"
        },
        "topLevelEnum": {
            "$ref": "samples.subpackage.EnumNumericValuesTopLevelEnum.json#"
        }
    },
    "additionalProperties": false,
    "type": "object",
    "definitions": {
        "EnumNumericValuesMsg.InlineEnum": {
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

const EnumNumericValuesMsg2 = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "additionalProperties": false,
    "type": "object",
    "definitions": {
        "EnumNumericValuesMsg2.NestedEnum": {
            "$schema": "http://json-schema.org/draft-04/schema#",
            "enum": [
                "FOO",
                0,
                "BAR",
                1,
                "FIZZ",
                2,
                "BUZZ",
                3,
                "BIZZ",
                4
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

const EnumNumericValuesTopLevelEnum = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "enum": [
        "FRR",
        0,
        "FRA",
        1
    ],
    "oneOf": [
        {
            "type": "string"
        },
        {
            "type": "integer"
        }
    ]
}`
