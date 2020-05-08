package testdata

const EnumNoNumericValuesMsg = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "properties": {
        "description": {
            "type": "string"
        },
        "stuff": {
            "$ref": "#/definitions/EnumNoNumericValuesMsg.InlineEnum"
        },
        "nestedEnum": {
            "$ref": "samples.subpackage.EnumNoNumericValuesMsg2.json#/definitions/EnumNoNumericValuesMsg2.NestedEnum"
        },
        "topLevelEnum": {
            "$ref": "samples.subpackage.EnumNoNumericValuesTopLevelEnum.json#"
        }
    },
    "additionalProperties": true,
    "type": "object",
    "definitions": {
        "EnumNoNumericValuesMsg.InlineEnum": {
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

const EnumNoNumericValuesMsg2 = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "additionalProperties": true,
    "type": "object",
    "definitions": {
        "EnumNoNumericValuesMsg2.NestedEnum": {
            "$schema": "http://json-schema.org/draft-04/schema#",
            "enum": [
                "FOO",
                "BAR",
                "FIZZ",
                "BUZZ",
                "BIZZ"
            ],
            "type": "string"
        }
    }
}`

const EnumNoNumericValuesTopLevelEnum = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "enum": [
        "FRR",
        "FRA"
    ],
    "type": "string"
}`
