package converter

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"sort"
	"strings"
)

type generatorPlan struct {
	targetFileLookup map[string][]*protoTypeInfo
	typeLookup       map[string]*protoTypeInfo
}

type protoTypeInfo struct {
	targetFileName        string
	jsonSchemaTopLevel    bool
	uniqueOnlyWithPackage bool
	protoPackage          []string
	protoName             []string
	protoMsg              *descriptor.DescriptorProto
	protoEnum             *descriptor.EnumDescriptorProto
}

func NewGeneratorPlan() *generatorPlan {
	plan := new(generatorPlan)
	plan.targetFileLookup = make(map[string][]*protoTypeInfo)
	plan.typeLookup = make(map[string]*protoTypeInfo)
	return plan
}

func (g *generatorPlan) Put(tInfo *protoTypeInfo) error {
	targetFileName := tInfo.GetTargetFileName()
	g.targetFileLookup[targetFileName] = append(g.targetFileLookup[targetFileName], tInfo)

	// Check if two types exists that have the same name but are located in different packages.
	// Normally the type name that is used in the OpenAPI schema does not contain the ProtoBuf package.
	// That makes it easier to read the OpenAPI specification. However, if to types have the same name the package
	// has to be appended to make them unique.
	for _, element := range g.typeLookup {
		if element.GetProtoTypeName() == tInfo.GetProtoTypeName() {
			element.uniqueOnlyWithPackage = true
			tInfo.uniqueOnlyWithPackage = true
		}
	}

	typeName := strings.Join(tInfo.getFullNameHierarchy(), ".")
	if g.typeLookup[typeName] != nil {
		return fmt.Errorf("type with full qualified name already exists: %s", typeName)
	}
	g.typeLookup[typeName] = tInfo

	return nil
}

func (g *generatorPlan) GetAllForTargetFilename(fileName string) []*protoTypeInfo {
	return g.targetFileLookup[fileName]
}

func (g *generatorPlan) GetAllTargetFilenames() []string {
	keys := make([]string, 0, len(g.targetFileLookup))
	for k, _ := range g.targetFileLookup {
		keys = append(keys, k)
	}

	// Sorting makes the order deterministic. Important for unit tests.
	sort.Strings(keys)
	return keys
}

func (g *generatorPlan) LookupType(typeName string) *protoTypeInfo {
	return g.typeLookup[typeName]
}

func (b *generatorPlan) String() string {
	result := ""
	for _, value := range b.targetFileLookup {
		for _, tInfo := range value {
			result += "[" + tInfo.String() + "]"
		}
	}
	return result
}

func NewProtoTypeInfoForMsg(targetFileName string, jsonSchemaTopLevel bool, protoPackage string, parent *protoTypeInfo,
	protoMsg *descriptor.DescriptorProto) *protoTypeInfo {

	return newProtoTypeInfo(targetFileName, protoPackage, parent, jsonSchemaTopLevel, nil, protoMsg)
}

func NewProtoTypeInfoForEnum(targetFileName string, jsonSchemaTopLevel bool, protoPackage string, parent *protoTypeInfo,
	protoEnum *descriptor.EnumDescriptorProto) *protoTypeInfo {

	return newProtoTypeInfo(targetFileName, protoPackage, parent, jsonSchemaTopLevel, protoEnum, nil)
}

func newProtoTypeInfo(targetFileName string,
	protoPackage string,
	parent *protoTypeInfo,
	jsonSchemaTopLevel bool,
	protoEnum *descriptor.EnumDescriptorProto,
	protoMsg *descriptor.DescriptorProto) *protoTypeInfo {

	tInfo := new(protoTypeInfo)
	tInfo.targetFileName = targetFileName
	tInfo.jsonSchemaTopLevel = jsonSchemaTopLevel
	tInfo.protoMsg = protoMsg
	tInfo.protoEnum = protoEnum
	tInfo.uniqueOnlyWithPackage = false

	tInfo.protoPackage = strings.Split(protoPackage, ".")
	if len(protoPackage) > 0 {
		// ProtoBuf implicitly uses a package with an empty name (i.e. "") as the first outermost package.
		// This creates a leading dot (.) in the full qualified type name. For example:
		// mypackage.MyType becomes .mypackage.MyType
		// MyType (no package) becomes .MyType
		// See https://developers.google.com/protocol-buffers/docs/proto3#packages-and-name-resolution
		tInfo.protoPackage = append([]string{""}, tInfo.protoPackage...)
	}

	tInfo.protoName = make([]string, 0, 10)
	if parent != nil {
		tInfo.protoName = append(tInfo.protoName, parent.protoName...)
	}

	if tInfo.protoMsg != nil {
		tInfo.protoName = append(tInfo.protoName, tInfo.protoMsg.GetName())
	} else {
		tInfo.protoName = append(tInfo.protoName, tInfo.protoEnum.GetName())
	}

	return tInfo
}

func (p *protoTypeInfo) GetTargetFileName() string {
	return p.targetFileName
}

func (p *protoTypeInfo) GenerateAtTopLevel() bool {
	return p.jsonSchemaTopLevel
}

func (p *protoTypeInfo) GetProtoTypeName() string {
	return strings.Trim(strings.Join(p.protoName, "."), ".")
}

func (p *protoTypeInfo) GetProtoFQTypeName() string {
	var name []string
	if p.uniqueOnlyWithPackage {
		name = append(p.protoPackage, p.protoName...)
	} else {
		name = p.protoName
	}

	return strings.Trim(strings.Join(name, "."), ".")
}

func (p *protoTypeInfo) IsProtoMessage() bool {
	return p.protoMsg != nil
}

func (p *protoTypeInfo) IsProtoEnum() bool {
	return p.protoEnum != nil
}

func (p *protoTypeInfo) GetProtoMessage() *descriptor.DescriptorProto {
	return p.protoMsg
}

func (p *protoTypeInfo) GetProtoEnum() *descriptor.EnumDescriptorProto {
	return p.protoEnum
}

func (p *protoTypeInfo) String() string {
	return p.targetFileName + " " + p.GetProtoFQTypeName()
}

func (p *protoTypeInfo) getFullNameHierarchy() []string {
	return append(p.protoPackage, p.protoName...)
}
