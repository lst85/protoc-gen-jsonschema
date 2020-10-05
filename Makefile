default: build-default

protoc-gen-%:
	@echo "Generating binary ($@) for os $(GOOS) and arch $(GOARCH) ..."
	@mkdir -p bin
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/$@ cmd/protoc-gen-openapi/main.go

build-all: build-linux build-windows build-darwin build-default

build-linux: GOOS = linux
build-linux: GOARCH = amd64
build-linux: protoc-gen-openapi.linux-amd64

build-windows: GOOS = windows
build-windows: GOARCH = amd64
build-windows: protoc-gen-openapi.windows-amd64

build-darwin: GOOS = darwin
build-darwin: GOARCH = amd64
build-darwin: protoc-gen-openapi.darwin-amd64

build-default: protoc-gen-openapi

proto_path ?= "internal/converter/testdata/proto"
comma:= ,

define build_sample
	@echo "Processing $(1) with parameters: $(2)"
	@PATH=./bin:$$PATH; protoc --openapi_out=$(2):jsonschemas --proto_path=$(proto_path) $(proto_path)/$(1) || echo "No messages found ($(1))"
endef

samples:
	@echo "Generating sample JSON-Schemas ..."
	@rm -f jsonschemas/*
	@mkdir -p jsonschemas
	$(call build_sample,AdditionalProperties.proto,out_file=samples.AdditionalProperties.json$(comma)allow_additional_properties)
	$(call build_sample,ArrayOfEnums.proto,out_file=samples.ArrayOfEnums.json)
	$(call build_sample,ArrayOfMessages.proto,out_file=samples.ArrayOfMessages.json)
	$(call build_sample,ArrayOfObjects.proto,out_file=samples.ArrayOfObjects.json$(comma)allow_null_values)
	$(call build_sample,ArrayOfPrimitives.proto,out_file=samples.ArrayOfPrimitives.json$(comma)allow_null_values)
	$(call build_sample,CyclicalReference.proto,out_file=samples.CyclicalReference.json)
	$(call build_sample,Enumception.proto,out_file=samples.Enumception.json$(comma)allow_additional_properties)
	$(call build_sample,EnumNumericValues.proto,out_file=samples.EnumNumericValues.json$(comma)allow_numeric_enum_values)
	$(call build_sample,Maps.proto,out_file=samples.Maps.json)
	$(call build_sample,MessageWithComments.proto,out_file=samples.MessageWithComments.json)
	$(call build_sample,NestedObject.proto,out_file=samples.NestedObject.json$(comma)disallow_bigints_as_strings)
	$(call build_sample,NoPackage.proto,out_file=samples.NoPackage.json)
	$(call build_sample,OpenApi.proto,out_file=openapi.json$(comma)open_api_template=${proto_path}/openapi.json)
	$(call build_sample,SelfReference.proto,out_file=samples.SelfReference.json$(comma)disallow_bigints_as_strings)
	$(call build_sample,SeveralEnums.proto,out_file=samples.SeveralEnums.json$(comma)disallow_bigints_as_strings)
	$(call build_sample,SeveralMessages.proto,out_file=samples.SeveralMessages.json)

protoc_zip := protoc-3.13.0-linux-x86_64.zip
download-protoc:
	@echo "Downloading $(protoc_zip)"
	@curl -OL -# "https://github.com/protocolbuffers/protobuf/releases/download/v3.13.0/$(protoc_zip)"
	@unzip -o "$(protoc_zip)" bin/protoc
	@rm -f "$(protoc_zip)"

test:
	@go test ./... -cover -v
