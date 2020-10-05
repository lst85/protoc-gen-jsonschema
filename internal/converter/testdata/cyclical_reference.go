package testdata

const (
	CyclicalReference = `{
    "components": {
        "schemas": {
            "Bar": {
                "properties": {
                    "id": {
                        "type": "integer"
                    },
                    "baz": {
                        "$ref": "#/components/schemas/Baz"
                    }
                },
                "type": "object"
            },
            "Baz": {
                "properties": {
                    "enabled": {
                        "type": "boolean"
                    },
                    "foo": {
                        "$ref": "#/components/schemas/Foo"
                    }
                },
                "type": "object"
            },
            "Foo": {
                "properties": {
                    "name": {
                        "type": "string"
                    },
                    "bar": {
                        "items": {
                            "$ref": "#/components/schemas/Bar"
                        },
                        "type": "array"
                    }
                },
                "type": "object"
            },
            "M": {
                "properties": {
                    "foo": {
                        "$ref": "#/components/schemas/Foo"
                    }
                },
                "type": "object"
            }
        }
    },
    "openapi": "3.0.0",
    "paths": {}
}`
)
