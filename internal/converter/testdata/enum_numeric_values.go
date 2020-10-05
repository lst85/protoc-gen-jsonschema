package testdata

const EnumNumericValues = `{
    "components": {
        "schemas": {
            "EnumNumericValuesMsg": {
                "properties": {
                    "description": {
                        "type": "string"
                    },
                    "stuff": {
                        "$ref": "#/components/schemas/EnumNumericValuesMsg.InlineEnum"
                    },
                    "nestedEnum": {
                        "$ref": "#/components/schemas/EnumNumericValuesMsg2.NestedEnum"
                    },
                    "topLevelEnum": {
                        "$ref": "#/components/schemas/EnumNumericValuesTopLevelEnum"
                    }
                },
                "type": "object"
            },
            "EnumNumericValuesMsg.InlineEnum": {
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
            },
            "EnumNumericValuesMsg2": {
                "type": "object"
            },
            "EnumNumericValuesMsg2.NestedEnum": {
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
            },
            "EnumNumericValuesTopLevelEnum": {
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
            }
        }
    },
    "openapi": "3.0.0",
    "paths": {}
}`
