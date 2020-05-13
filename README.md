![Go](https://github.com/lst85/protoc-gen-jsonschema/workflows/Go/badge.svg?branch=master)

Protobuf to JSON Schema and OpenAPI v3 compiler
================================================
This is a protoc plugin that takes protocol buffers definitions and converts them into JSON Schemas or OpenAPI v.3 documents.

This will hopefully be useful for people who define their data using ProtoBuf, but use JSON for the "wire" format.

Forked from [chrusty/protoc-gen-jsonschema](https://github.com/chrusty/protoc-gen-jsonschema) and
"heavily influenced" by [Google's protobuf-to-BigQuery-schema compiler](https://github.com/GoogleCloudPlatform/protoc-gen-bq-schema).


Installation
------------

First install Go and [protoc](https://github.com/protocolbuffers/protobuf), then install the plugin with:

`GO111MODULE=on go get github.com/lst85/protoc-gen-jsonschema/cmd/protoc-gen-jsonschema && go install github.com/lst85/protoc-gen-jsonschema/cmd/protoc-gen-jsonschema`

Links
-----
* [About JSON Schema](http://json-schema.org/)
* [Proto3 Language Guide](https://developers.google.com/protocol-buffers/docs/proto3#json)
* [OpenAPI Specification Version 3.0.0](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.0.md)

Usage
-----

Simply invoke `protoc` with the `--jsonschema_out` command-line parameter:

```
protoc --jsonschema_out="<OUT_DIR>" input.proto
protoc --jsonschema_out="<OPTIONS>:<OUT_DIR>" input.proto
```

Where `<OPTIONS>` is a comma-seperated list of `key=value,key2` options (which are listed in detail below). 

For example:

`protoc --jsonschema_out="my_schemas/" file.proto`
would produce JSON schemas for `file.proto` inside the `my_schemas/` directory. 

`protoc --jsonschema_out="open_api,open_api_template=template.json:my_schemas/" file.proto`
would produce an OpenAPI document for `file.proto` inside the `my_schemas/` directory using the template from `template.json`.

Options
-----

| Option              | Description |
|---------------------|-------------|
| `allow_null_values` | Allow NULL values for all properties. By default, JSON Schemas will reject NULL values. |
| `allow_additional_properties` | Allow additional properties. JSON Schemas will allow extra parameters, that are not specified in the schema. |
| `debug` | Enable debug logging. |
| `disallow_bigints_as_strings` | If the parameter is not set (default) the JSON Schema will allow both string and integers for 64 bit integers. If it is set only integers are allowed. The canonical JSON encoding of Proto3 converts int64, fixed64, uint6 to JSON strings. When decoding JSON to ProtoBuf both numbers and strings are accepted. |
| `allow_numeric_enum_values` | Allow both enum names and integer values. |
| `out_file=<file>` | Create a single file instead of multiple files. When JSON Schema mode is enabled (parameter `open_api` is not set) a single JSON Schema will be generated with the given filename. When OpenAPI mode is enabled (parameter `open_api` is set) this parameter is implicitly enabled and the default filename is `openapi.json`. |
| `open_api` | Generate an OpenAPI v.3 file instead of JSON Schema file(s). All ProtoBuf types (messages, enums, etc.) will be converted to their JSON Schema equivalent and added to the components/schemas section of the OpenAPI document. *The generator currently ignores gRPC service definitions.* The paths section of the generated OpenAPI document will be emtpy. |
| `open_api_template=<file>` | Path to an OpenAPI file that will be merged with the generated schemas. This parameter has only an effect when the parameter `open_api` is set. |
| `proto_fieldnames` | If the parameter is set the field names from the ProtoBuf definition are used in the JSON Schema. If the parameter is not set message field names are mapped to lowerCamelCase and become JSON object keys. |