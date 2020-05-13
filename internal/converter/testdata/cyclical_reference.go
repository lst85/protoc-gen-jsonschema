package testdata

const (
	CyclicalReferenceMessageM = `{
    "$schema": "http://json-schema.org/draft-04/schema#",
    "properties": {
        "foo": {
            "$ref": "samples.cyclicalreference.Foo.json#",
            "additionalProperties": false,
            "type": "object"
        }
    },
    "additionalProperties": false,
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
    "additionalProperties": false,
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
            "additionalProperties": false,
            "type": "object"
        }
    },
    "additionalProperties": false,
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
            "additionalProperties": false,
            "type": "object"
        }
    },
    "additionalProperties": false,
    "type": "object"
}`
)
