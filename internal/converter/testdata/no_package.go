package testdata

const NoPackageMsg = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "properties": {
        "payload": {
            "$ref": "#/definitions/NoPackageMsg.NestedMsg",
            "additionalProperties": false
        },
        "description": {
            "type": "string"
        },
        "nestEnum": {
            "$ref": "#/definitions/NoPackageMsg.NestedMsg.NestedEnum"
        },
        "otherMsg": {
            "$ref": "NoPackageOtherMsg.json#",
            "additionalProperties": false
        }
    },
    "additionalProperties": false,
    "type": "object",
    "definitions": {
        "NoPackageMsg.NestedMsg": {
            "$schema": "http://json-schema.org/draft-04/schema#",
            "properties": {
                "name": {
                    "type": "string"
                },
                "nestEnum": {
                    "$ref": "#/definitions/NoPackageMsg.NestedMsg.NestedEnum"
                }
            },
            "additionalProperties": false,
            "type": "object"
        },
        "NoPackageMsg.NestedMsg.NestedEnum": {
            "$schema": "http://json-schema.org/draft-04/schema#",
            "enum": [
                "FLAT",
                "NESTED_OBJECT"
            ],
            "type": "string"
        }
    }
}`

const NoPackageOtherMsg = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "properties": {
        "prop": {
            "type": "string"
        },
        "msg": {
            "$ref": "NoPackageMsg.json#",
            "additionalProperties": false
        },
        "nestEnum": {
            "$ref": "NoPackageMsg.json#/definitions/NoPackageMsg.NestedMsg.NestedEnum"
        }
    },
    "additionalProperties": false,
    "type": "object"
}`
