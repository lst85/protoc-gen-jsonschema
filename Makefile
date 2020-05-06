default: build-default

protoc-gen-%:
	@echo "Generating binary ($@) for os $(GOOS) and arch $(GOARCH) ..."
	@mkdir -p bin
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/$@ cmd/protoc-gen-jsonschema/main.go

build-all: build-linux build-windows build-darwin build-default

build-linux: GOOS = linux
build-linux: GOARCH = amd64
build-linux: protoc-gen-jsonschema.linux-amd64

build-windows: GOOS = linux
build-windows: GOARCH = amd64
build-windows: protoc-gen-jsonschema.windows-amd64

build-darwin: GOOS = linux
build-darwin: GOARCH = amd64
build-darwin: protoc-gen-jsonschema.darwin-amd64

build-default: protoc-gen-jsonschema

PROTO_PATH ?= "internal/converter/testdata/proto"
samples:
	@echo "Generating sample JSON-Schemas ..."
	@mkdir -p jsonschemas
	@PATH=./bin:$$PATH; protoc --jsonschema_out=open_api_conform:jsonschemas --proto_path=${PROTO_PATH} ${PROTO_PATH}/AdditionalProperties.proto || echo "No messages found (AdditionalProperties.proto)"
	@PATH=./bin:$$PATH; protoc --jsonschema_out=allow_null_values:jsonschemas --proto_path=${PROTO_PATH} ${PROTO_PATH}/ArrayOfMessages.proto || echo "No messages found (ArrayOfMessages.proto)"
	@PATH=./bin:$$PATH; protoc --jsonschema_out=allow_null_values:jsonschemas --proto_path=${PROTO_PATH} ${PROTO_PATH}/ArrayOfObjects.proto || echo "No messages found (ArrayOfObjects.proto)"
	@PATH=./bin:$$PATH; protoc --jsonschema_out=allow_null_values:jsonschemas --proto_path=${PROTO_PATH} ${PROTO_PATH}/ArrayOfPrimitives.proto || echo "No messages found (ArrayOfPrimitives.proto)"
	@PATH=./bin:$$PATH; protoc --jsonschema_out=disallow_additional_properties:jsonschemas --proto_path=${PROTO_PATH} ${PROTO_PATH}/Enumception.proto || echo "No messages found (Enumception.proto)"
	@PATH=./bin:$$PATH; protoc --jsonschema_out=disallow_numeric_enum_values:jsonschemas --proto_path=${PROTO_PATH} ${PROTO_PATH}/EnumNoNumericValues.proto || echo "No messages found (EnumNoNumericValues.proto)"
	@PATH=./bin:$$PATH; protoc --jsonschema_out=disallow_additional_properties:jsonschemas --proto_path=${PROTO_PATH} ${PROTO_PATH}/ImportedEnum.proto || echo "No messages found (ImportedEnum.proto)"
	@PATH=./bin:$$PATH; protoc --jsonschema_out=disallow_additional_properties:jsonschemas --proto_path=${PROTO_PATH} ${PROTO_PATH}/NestedMessage.proto || echo "No messages found (NestedMessage.proto)"
	@PATH=./bin:$$PATH; protoc --jsonschema_out=disallow_bigints_as_strings:jsonschemas --proto_path=${PROTO_PATH} ${PROTO_PATH}/NestedObject.proto || echo "No messages found (NestedObject.proto)"
	@PATH=./bin:$$PATH; protoc --jsonschema_out=disallow_bigints_as_strings:jsonschemas --proto_path=${PROTO_PATH} ${PROTO_PATH}/PayloadMessage.proto || echo "No messages found (PayloadMessage.proto)"
	@PATH=./bin:$$PATH; protoc --jsonschema_out=disallow_bigints_as_strings:jsonschemas --proto_path=${PROTO_PATH} ${PROTO_PATH}/SeveralEnums.proto || echo "No messages found (SeveralEnums.proto)"
	@PATH=./bin:$$PATH; protoc --jsonschema_out=disallow_bigints_as_strings:jsonschemas --proto_path=${PROTO_PATH} ${PROTO_PATH}/SeveralMessages.proto || echo "No messages found (SeveralMessages.proto)"
	@PATH=./bin:$$PATH; protoc --jsonschema_out=jsonschemas --proto_path=${PROTO_PATH} ${PROTO_PATH}/ArrayOfEnums.proto || echo "No messages found (SeveralMessages.proto)"
	@PATH=./bin:$$PATH; protoc --jsonschema_out=jsonschemas --proto_path=${PROTO_PATH} ${PROTO_PATH}/Maps.proto || echo "No messages found (Maps.proto)"
	@PATH=./bin:$$PATH; protoc --jsonschema_out=jsonschemas --proto_path=${PROTO_PATH} ${PROTO_PATH}/MessageWithComments.proto || echo "No messages found (MessageWithComments.proto)"

test:
	@go test ./... -cover -v
