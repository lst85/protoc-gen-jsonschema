package testdata

const NoPackage = `{
    "components": {
        "schemas": {
            "NoPackageMsg": {
                "properties": {
                    "payload": {
                        "$ref": "#/components/schemas/NoPackageMsg.NestedMsg"
                    },
                    "description": {
                        "type": "string"
                    },
                    "nestEnum": {
                        "$ref": "#/components/schemas/NoPackageMsg.NestedMsg.NestedEnum"
                    },
                    "otherMsg": {
                        "$ref": "#/components/schemas/NoPackageOtherMsg"
                    }
                },
                "type": "object"
            },
            "NoPackageMsg.NestedMsg": {
                "properties": {
                    "name": {
                        "type": "string"
                    },
                    "nestEnum": {
                        "$ref": "#/components/schemas/NoPackageMsg.NestedMsg.NestedEnum"
                    }
                },
                "type": "object"
            },
            "NoPackageMsg.NestedMsg.NestedEnum": {
                "enum": [
                    "FLAT",
                    "NESTED_OBJECT"
                ],
                "type": "string"
            },
            "NoPackageOtherMsg": {
                "properties": {
                    "prop": {
                        "type": "string"
                    },
                    "msg": {
                        "$ref": "#/components/schemas/NoPackageMsg"
                    },
                    "nestEnum": {
                        "$ref": "#/components/schemas/NoPackageMsg.NestedMsg.NestedEnum"
                    }
                },
                "type": "object"
            }
        }
    },
    "openapi": "3.0.0",
    "paths": {}
}`
