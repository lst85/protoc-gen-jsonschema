package converter

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/chrusty/protoc-gen-jsonschema/internal/converter/testdata"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	sampleProtoDirectory = "testdata/proto"
	sampleProtos         = make(map[string]sampleProto)
)

type sampleProto struct {
	AllowNullValues           bool
	AllowAdditionalProperties bool
	AllowNumericEnumValues    bool
	DisallowBigIntsAsStrings  bool
	GenerateOpenApi           bool
	OpenApiFile               string
	ExpectedJSONSchema        []string
	FilesToGenerate           []string
	ProtoFileName             string
	UseProtoAndJSONFieldNames bool
}

func TestGenerateJsonSchema(t *testing.T) {

	// Configure the list of sample protos to test, and their expected JSON-Schemas:
	configureSampleProtos()

	// Convert the protos, compare the results against the expected JSON-Schemas:
	testConvertSampleProto(t, sampleProtos["AdditionalProperties"])
	testConvertSampleProto(t, sampleProtos["ArrayOfEnums"])
	testConvertSampleProto(t, sampleProtos["ArrayOfMessages"])
	testConvertSampleProto(t, sampleProtos["ArrayOfObjects"])
	testConvertSampleProto(t, sampleProtos["ArrayOfPrimitives"])
	testConvertSampleProto(t, sampleProtos["CyclicalReference"])
	testConvertSampleProto(t, sampleProtos["Enumception"])
	testConvertSampleProto(t, sampleProtos["EnumNumericValues"])
	testConvertSampleProto(t, sampleProtos["Maps"])
	testConvertSampleProto(t, sampleProtos["MessageWithComments"])
	testConvertSampleProto(t, sampleProtos["NestedObject"])
	testConvertSampleProto(t, sampleProtos["NoPackage"])
	testConvertSampleProto(t, sampleProtos["OpenApi"])
	testConvertSampleProto(t, sampleProtos["SelfReference"])
	testConvertSampleProto(t, sampleProtos["SeveralEnums"])
	testConvertSampleProto(t, sampleProtos["SeveralMessages"])
}

func testConvertSampleProto(t *testing.T, sampleProto sampleProto) {

	// Make a Logrus logger:
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	logger.SetOutput(os.Stderr)

	// Use the logger to make a Converter:
	protoConverter := New(logger)
	protoConverter.AllowNullValues = sampleProto.AllowNullValues
	protoConverter.AllowNumericEnumValues = sampleProto.AllowNumericEnumValues
	protoConverter.AllowAdditionalProperties = sampleProto.AllowAdditionalProperties
	protoConverter.DisallowBigIntsAsStrings = sampleProto.DisallowBigIntsAsStrings
	protoConverter.GenerateOpenApi = sampleProto.GenerateOpenApi
	if protoConverter.GenerateOpenApi {
		protoConverter.SingleOutputFile = sampleProto.ProtoFileName + "_openapi.json"
	}
	protoConverter.OpenApiFile = sampleProto.OpenApiFile
	protoConverter.UseProtoFieldNames = sampleProto.UseProtoAndJSONFieldNames

	// Open the sample proto file:
	sampleProtoFileName := fmt.Sprintf("%v/%v", sampleProtoDirectory, sampleProto.ProtoFileName)
	fileDescriptorSet := mustReadProtoFiles(t, sampleProtoDirectory, sampleProto.ProtoFileName)

	// Prepare a request:
	codeGeneratorRequest := plugin.CodeGeneratorRequest{
		FileToGenerate: sampleProto.FilesToGenerate,
		ProtoFile:      fileDescriptorSet.GetFile(),
	}

	// Perform the conversion:
	response, err := protoConverter.convert(&codeGeneratorRequest)
	assert.NoError(t, err, "Unable to convert sample proto file (%v)", sampleProtoFileName)
	assert.Equal(t, len(sampleProto.ExpectedJSONSchema), len(response.File),
		"Incorrect number of JSON-Schema files returned for sample proto file (%v)", sampleProtoFileName)
	if len(sampleProto.ExpectedJSONSchema) != len(response.File) {
		t.Fail()
	} else {
		for responseFileIndex, responseFile := range response.File {
			assert.Equal(t, sampleProto.ExpectedJSONSchema[responseFileIndex], *responseFile.Content,
				"Incorrect JSON-Schema returned for sample proto file (%v)", sampleProtoFileName)
		}
	}

}

func configureSampleProtos() {
	// AdditionalProperties:
	sampleProtos["AdditionalProperties"] = sampleProto{
		GenerateOpenApi:           true,
		AllowAdditionalProperties: true,
		ExpectedJSONSchema:        []string{testdata.AdditionalProperties},
		FilesToGenerate:           []string{"AdditionalProperties.proto"},
		ProtoFileName:             "AdditionalProperties.proto",
	}

	// ArrayOfEnums:
	sampleProtos["ArrayOfEnums"] = sampleProto{
		ExpectedJSONSchema: []string{testdata.ArrayOfEnums},
		FilesToGenerate:    []string{"ArrayOfEnums.proto"},
		ProtoFileName:      "ArrayOfEnums.proto",
	}

	// ArrayOfMessages:
	sampleProtos["ArrayOfMessages"] = sampleProto{
		ExpectedJSONSchema: []string{testdata.ArrayOfMessages},
		FilesToGenerate:    []string{"ArrayOfMessages.proto"},
		ProtoFileName:      "ArrayOfMessages.proto",
	}

	// ArrayOfObjects:
	sampleProtos["ArrayOfObjects"] = sampleProto{
		AllowNullValues:    true,
		ExpectedJSONSchema: []string{testdata.ArrayOfObjects},
		FilesToGenerate:    []string{"ArrayOfObjects.proto"},
		ProtoFileName:      "ArrayOfObjects.proto",
	}

	// ArrayOfPrimitives:
	sampleProtos["ArrayOfPrimitives"] = sampleProto{
		AllowNullValues:    true,
		ExpectedJSONSchema: []string{testdata.ArrayOfPrimitives},
		FilesToGenerate:    []string{"ArrayOfPrimitives.proto"},
		ProtoFileName:      "ArrayOfPrimitives.proto",
	}

	// Messages that depend on one another so as to form a cycle
	sampleProtos["CyclicalReference"] = sampleProto{
		ExpectedJSONSchema: []string{testdata.CyclicalReferenceMessageBar, testdata.CyclicalReferenceMessageBaz,
			testdata.CyclicalReferenceMessageFoo, testdata.CyclicalReferenceMessageM},
		FilesToGenerate: []string{"CyclicalReference.proto"},
		ProtoFileName:   "CyclicalReference.proto",
	}

	// EnumCeption:
	sampleProtos["Enumception"] = sampleProto{
		AllowAdditionalProperties: true,
		ExpectedJSONSchema:        []string{testdata.Enumception, testdata.EnumceptionImportedEnum, testdata.EnumceptionImportedMessage},
		FilesToGenerate:           []string{"Enumception.proto", "_ImportedMessage.proto", "_ImportedEnum.proto"},
		ProtoFileName:             "Enumception.proto",
	}

	// Tests the DisallowNumericEnumValues parameter:
	sampleProtos["EnumNumericValues"] = sampleProto{
		AllowNumericEnumValues: true,
		ExpectedJSONSchema: []string{testdata.EnumNumericValuesMsg, testdata.EnumNumericValuesMsg2,
			testdata.EnumNumericValuesTopLevelEnum},
		FilesToGenerate: []string{"EnumNumericValues.proto"},
		ProtoFileName:   "EnumNumericValues.proto",
	}

	// Maps:
	sampleProtos["Maps"] = sampleProto{
		ExpectedJSONSchema: []string{testdata.Maps},
		FilesToGenerate:    []string{"Maps.proto"},
		ProtoFileName:      "Maps.proto",
	}

	// MessageWithComments:
	sampleProtos["MessageWithComments"] = sampleProto{
		ExpectedJSONSchema: []string{testdata.MessageWithComments},
		FilesToGenerate:    []string{"MessageWithComments.proto"},
		ProtoFileName:      "MessageWithComments.proto",
	}

	// NestedObject:
	sampleProtos["NestedObject"] = sampleProto{
		DisallowBigIntsAsStrings: true,
		ExpectedJSONSchema:       []string{testdata.NestedObject},
		FilesToGenerate:          []string{"NestedObject.proto"},
		ProtoFileName:            "NestedObject.proto",
	}

	// NoPackage:
	sampleProtos["NoPackage"] = sampleProto{
		DisallowBigIntsAsStrings: true,
		ExpectedJSONSchema:       []string{testdata.NoPackageMsg, testdata.NoPackageOtherMsg},
		FilesToGenerate:          []string{"NoPackage.proto"},
		ProtoFileName:            "NoPackage.proto",
	}

	// OpenApi:
	sampleProtos["OpenApi"] = sampleProto{
		GenerateOpenApi:    true,
		OpenApiFile:        fmt.Sprintf("%v/openapi.json", sampleProtoDirectory),
		ExpectedJSONSchema: []string{testdata.OpenApi},
		FilesToGenerate:    []string{"OpenApi.proto"},
		ProtoFileName:      "OpenApi.proto",
	}

	// Self referencing proto message
	sampleProtos["SelfReference"] = sampleProto{
		DisallowBigIntsAsStrings: true,
		ExpectedJSONSchema:       []string{testdata.SelfReference},
		FilesToGenerate:          []string{"SelfReference.proto"},
		ProtoFileName:            "SelfReference.proto",
	}

	// SeveralEnums:
	sampleProtos["SeveralEnums"] = sampleProto{
		DisallowBigIntsAsStrings: true,
		ExpectedJSONSchema:       []string{testdata.SeveralEnumsFirstEnum, testdata.SeveralEnumsSecondEnum},
		FilesToGenerate:          []string{"SeveralEnums.proto"},
		ProtoFileName:            "SeveralEnums.proto",
	}

	// SeveralMessages:
	sampleProtos["SeveralMessages"] = sampleProto{
		ExpectedJSONSchema: []string{testdata.FirstMessage, testdata.SecondMessage},
		FilesToGenerate:    []string{"SeveralMessages.proto"},
		ProtoFileName:      "SeveralMessages.proto",
	}
}

// Load the specified .proto files into a FileDescriptorSet. Any errors in loading/parsing will
// immediately fail the test.
func mustReadProtoFiles(t *testing.T, includePath string, filenames ...string) *descriptor.FileDescriptorSet {
	protocBinary, err := exec.LookPath("protoc")
	if err != nil {
		t.Fatalf("Can't find 'protoc' binary in $PATH: %s", err.Error())
	}

	// Use protoc to output descriptor info for the specified .proto files.
	var args []string
	args = append(args, "--descriptor_set_out=/dev/stdout")
	args = append(args, "--include_source_info")
	args = append(args, "--include_imports")
	args = append(args, "--proto_path="+includePath)
	args = append(args, filenames...)
	cmd := exec.Command(protocBinary, args...)
	stdoutBuf := bytes.Buffer{}
	stderrBuf := bytes.Buffer{}
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err = cmd.Run()
	if err != nil {
		t.Fatalf("failed to load descriptor set (%s): %s: %s",
			strings.Join(cmd.Args, " "), err.Error(), stderrBuf.String())
	}
	fds := &descriptor.FileDescriptorSet{}
	err = proto.Unmarshal(stdoutBuf.Bytes(), fds)
	if err != nil {
		t.Fatalf("failed to parse protoc output as FileDescriptorSet: %s", err.Error())
	}
	return fds
}
