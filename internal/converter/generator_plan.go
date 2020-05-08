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
	targetFileName     string
	jsonSchemaTopLevel bool
	protoPackage       []string
	protoFQName        []string
	protoMsg           *descriptor.DescriptorProto
	protoEnum          *descriptor.EnumDescriptorProto
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
	typeName = strings.Trim(typeName, ".")
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

	tInfo.protoPackage = strings.Split(protoPackage, ".")

	tInfo.protoFQName = make([]string, 0, 10)
	if parent != nil {
		tInfo.protoFQName = append(tInfo.protoFQName, parent.protoFQName...)
	}

	if tInfo.protoMsg != nil {
		tInfo.protoFQName = append(tInfo.protoFQName, tInfo.protoMsg.GetName())
	} else {
		tInfo.protoFQName = append(tInfo.protoFQName, tInfo.protoEnum.GetName())
	}

	return tInfo
}

func (p *protoTypeInfo) GetTargetFileName() string {
	return p.targetFileName
}

func (p *protoTypeInfo) GenerateAtTopLevel() bool {
	return p.jsonSchemaTopLevel
}

func (p *protoTypeInfo) GetProtoFQTypeName() string {
	return strings.Join(p.protoFQName, ".")
}

func (p *protoTypeInfo) GetProtoTypeName() string {
	return strings.Join(append(p.protoPackage, p.protoFQName...), ".")
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
	return append(p.protoPackage, p.protoFQName...)
}
