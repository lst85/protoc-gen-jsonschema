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

proto_path ?= "internal/converter/testdata/proto"
comma:= ,

define build_sample
	@echo "Processing $(1) with parameters: $(2)"
	@PATH=./bin:$$PATH; protoc --jsonschema_out=$(2):jsonschemas --proto_path=$(proto_path) $(proto_path)/$(1) || echo "No messages found ($(1))"
endef

samples:
	@echo "Generating sample JSON-Schemas ..."
	@rm -f jsonschemas/*
	@mkdir -p jsonschemas
	$(call build_sample,AdditionalProperties.proto,open_api$(comma)out_file=samples.AdditionalProperties.json$(comma)allow_additional_properties)
	$(call build_sample,ArrayOfEnums.proto,)
	$(call build_sample,ArrayOfMessages.proto,)
	$(call build_sample,ArrayOfObjects.proto,allow_null_values)
	$(call build_sample,ArrayOfPrimitives.proto,allow_null_values)
	$(call build_sample,CyclicalReference.proto,)
	$(call build_sample,Enumception.proto,allow_additional_properties)
	$(call build_sample,EnumNumericValues.proto,allow_numeric_enum_values)
	$(call build_sample,Maps.proto,)
	$(call build_sample,MessageWithComments.proto,)
	$(call build_sample,NestedObject.proto,disallow_bigints_as_strings)
	$(call build_sample,NoPackage.proto,)
	$(call build_sample,OpenApi.proto,open_api$(comma)out_file=openapi.json$(comma)open_api_template=${proto_path}/openapi.json)
	$(call build_sample,SelfReference.proto,disallow_bigints_as_strings)
	$(call build_sample,SeveralEnums.proto,disallow_bigints_as_strings)
	$(call build_sample,SeveralMessages.proto,)

test:
	@go test ./... -cover -v
