package converter

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"strings"
)

type generatorPlan struct {
	jsonFileLookup map[string][]*protoTypeInfo
	typeLookup     map[string]*protoTypeInfo
}

type protoTypeInfo struct {
	jsonFileName string
	jsonTopLevel bool
	protoPackage []string
	protoFQName  []string
	protoMsg     *descriptor.DescriptorProto
	protoEnum    *descriptor.EnumDescriptorProto
	children     []*protoTypeInfo
}

func NewGeneratorPlan() *generatorPlan {
	plan := new(generatorPlan)
	plan.jsonFileLookup = make(map[string][]*protoTypeInfo)
	plan.typeLookup = make(map[string]*protoTypeInfo)
	return plan
}

func (g *generatorPlan) Put(tInfo *protoTypeInfo) error {
	jsonFileName := tInfo.GetJsonFileName()
	g.jsonFileLookup[jsonFileName] = append(g.jsonFileLookup[jsonFileName], tInfo)

	typeName := strings.Join(tInfo.getFullNameHierarchy(), ".")
	if g.typeLookup[typeName] != nil {
		return fmt.Errorf("type with full qualified name already exists: %s", typeName)
	}
	g.typeLookup[typeName] = tInfo
	return nil
}

func (g *generatorPlan) GetForJsonFilename(fileName string) []*protoTypeInfo {
	return g.jsonFileLookup[fileName]
}

func (g *generatorPlan) GetAllJsonFilenames() []string {
	keys := make([]string, 0, len(g.jsonFileLookup))
	for k, _ := range g.jsonFileLookup {
		keys = append(keys, k)
	}
	return keys
}

func (g *generatorPlan) LookupType(contextType *protoTypeInfo, typeName string) *protoTypeInfo {

	typeName = strings.Trim(typeName, ".")
	return g.typeLookup[typeName]
}

func (b *generatorPlan) String() string {
	result := ""
	for _, value := range b.jsonFileLookup {
		for _, tInfo := range value {
			result += "[" + tInfo.String() + "]"
		}
	}
	return result
}

func NewProtoTypeInfoForMsg(jsonFileName string, jsonTopLevel bool, protoPackage string, parent *protoTypeInfo,
	protoMsg *descriptor.DescriptorProto) *protoTypeInfo {

	return newProtoTypeInfo(jsonFileName, protoPackage, parent, jsonTopLevel, nil, protoMsg)
}

func NewProtoTypeInfoForEnum(jsonFileName string, jsonTopLevel bool, protoPackage string, parent *protoTypeInfo,
	protoEnum *descriptor.EnumDescriptorProto) *protoTypeInfo {

	return newProtoTypeInfo(jsonFileName, protoPackage, parent, jsonTopLevel, protoEnum, nil)
}

func newProtoTypeInfo(jsonFileName string,
	protoPackage string,
	parent *protoTypeInfo,
	jsonTopLevel bool,
	protoEnum *descriptor.EnumDescriptorProto,
	protoMsg *descriptor.DescriptorProto) *protoTypeInfo {

	tInfo := new(protoTypeInfo)
	tInfo.jsonFileName = jsonFileName
	tInfo.jsonTopLevel = jsonTopLevel
	tInfo.protoMsg = protoMsg
	tInfo.protoEnum = protoEnum
	tInfo.children = make([]*protoTypeInfo, 0, 10)

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

func (p *protoTypeInfo) GetJsonFileName() string {
	return p.jsonFileName
}

func (p *protoTypeInfo) GetJsonRef(contextType *protoTypeInfo) string {
	ref := ""
	if contextType.GetJsonFileName() != p.GetJsonFileName() {
		ref += p.GetJsonFileName()
	}
	ref += "#"
	if !p.GenerateAtTopLevel() {
		ref += "/definitions/"
		ref += p.GetProtoFQNName()
	}
	return ref
}

func (p *protoTypeInfo) GenerateAtTopLevel() bool {
	return p.jsonTopLevel
}

func (p *protoTypeInfo) GetProtoTypeName() string {
	return p.protoFQName[len(p.protoFQName)-1]
}

func (p *protoTypeInfo) GetProtoFQNName() string {
	return strings.Join(p.protoFQName, ".")
}

func (p *protoTypeInfo) GetProtoPackageName() string {
	return strings.Join(p.protoPackage, ".")
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
	return p.jsonFileName + " " + p.GetProtoPackageName() + " " + p.GetProtoFQNName()
}

func (p *protoTypeInfo) getFullNameHierarchy() []string {
	return append(p.protoPackage, p.protoFQName...)
}
