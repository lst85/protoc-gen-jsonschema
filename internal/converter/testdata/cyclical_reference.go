package testdata

const (
	CyclicalReferenceMessageM = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "properties": {
        "foo": {
            "$ref": "samples.cyclicalreference.Foo.json#",
            "additionalProperties": true,
            "type": "object"
        }
    },
    "additionalProperties": true,
    "type": "object"
}`

	CyclicalReferenceMessageFoo = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "properties": {
        "name": {
            "type": "string"
        },
        "bar": {
            "items": {
                "$ref": "samples.cyclicalreference.Bar.json#"
            },
            "type": "array"
        }
    },
    "additionalProperties": true,
    "type": "object"
}`

	CyclicalReferenceMessageBar = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "properties": {
        "id": {
            "type": "integer"
        },
        "baz": {
            "$ref": "samples.cyclicalreference.Baz.json#",
            "additionalProperties": true,
            "type": "object"
        }
    },
    "additionalProperties": true,
    "type": "object"
}`

	CyclicalReferenceMessageBaz = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "properties": {
        "enabled": {
            "type": "boolean"
        },
        "foo": {
            "$ref": "samples.cyclicalreference.Foo.json#",
            "additionalProperties": true,
            "type": "object"
        }
    },
    "additionalProperties": true,
    "type": "object"
}`
)
