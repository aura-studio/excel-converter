package converter

import (
	"fmt"
	"sort"
	"strings"
)

type Identifier struct {
	OriginNodes []Node

	StrNodeMap map[string]uint32 // key: str, val: node id
	NodeStrMap map[uint32]string // key: node id, val: str

	NodeStructMap   map[uint32]string // key: node id, val: struct name
	NodeTypeMap     map[uint32]string // key: node id, val: type name
	NodeDataTypeMap map[uint32]string // key: node id, val: type name only struct type

	StructEqualMap map[string]string // key: dst struct, val: src struct
	StructEquals   [][]string
}

func NewIdentifier() *Identifier {
	return &Identifier{
		StrNodeMap:      make(map[string]uint32),
		NodeStrMap:      make(map[uint32]string),
		NodeStructMap:   make(map[uint32]string),
		NodeTypeMap:     make(map[uint32]string),
		NodeDataTypeMap: make(map[uint32]string),
		StructEqualMap:  make(map[string]string),
	}
}

func (i *Identifier) GenerateStruct(node Node) {
	i.generateStruct(node.RootName(), node)
	sort.Slice(i.OriginNodes, func(m, n int) bool {
		return i.NodeStructMap[i.OriginNodes[m].ID()] < i.NodeStructMap[i.OriginNodes[n].ID()]
	})
}

func (i *Identifier) GenerateType(node Node) {
	i.generateType(node)
}

func (i *Identifier) generateStruct(name string, node Node) {
	name = fmt.Sprintf("%s%s", name, node.Field().Name)
	i.NodeStructMap[node.ID()] = name
	for _, node := range node.Nodes() {
		i.generateStruct(name, node)
	}
}

func (i *Identifier) generateType(node Node) {
	switch node.Type() {
	case NodeTypeStruct:
		for _, node := range node.Nodes() {
			i.generateType(node)
		}
		i.NodeTypeMap[node.ID()] = fmt.Sprintf("*%s", i.NodeStructMap[node.ID()])
		i.NodeDataTypeMap[node.ID()] = fmt.Sprintf("*structs.%s", i.NodeStructMap[node.ID()])
	case NodeTypeSlice:
		i.generateType(node.Nodes()[0])
		i.NodeTypeMap[node.ID()] = fmt.Sprintf("[]%s", i.NodeTypeMap[node.Nodes()[0].ID()])
		i.NodeDataTypeMap[node.ID()] = fmt.Sprintf("[]%s", i.NodeDataTypeMap[node.Nodes()[0].ID()])
	case NodeTypeMap:
		i.generateType(node.Nodes()[0])
		i.generateType(node.Nodes()[1])
		i.NodeTypeMap[node.ID()] = fmt.Sprintf("map[%s]%s",
			i.NodeTypeMap[node.Nodes()[0].ID()], i.NodeTypeMap[node.Nodes()[1].ID()])
		i.NodeDataTypeMap[node.ID()] = fmt.Sprintf("map[%s]%s",
			i.NodeDataTypeMap[node.Nodes()[0].ID()], i.NodeDataTypeMap[node.Nodes()[1].ID()])
	case NodeTypeSimple:
		switch node.Field().Structure {
		case StructureTypeString:
			i.NodeTypeMap[node.ID()] = "string"
		case StructureTypeInt:
			i.NodeTypeMap[node.ID()] = "int64"
		case StructureTypeFloat:
			i.NodeTypeMap[node.ID()] = "float64"
		case StructureTypeBool:
			i.NodeTypeMap[node.ID()] = "bool"
		case StructureTypeBigInt:
			i.NodeTypeMap[node.ID()] = "*big.Int"
		case StructureTypeBigFloat:
			i.NodeTypeMap[node.ID()] = "*big.Float"
		case StructureTypeBigRat:
			i.NodeTypeMap[node.ID()] = "*big.Rat"
		case StructureTypeTime:
			i.NodeTypeMap[node.ID()] = "*time.Time"
		}
		i.NodeDataTypeMap[node.ID()] = i.NodeTypeMap[node.ID()]
	default:
		Exit("[%v] Unsupported node type %v", node.Sheet(), node.Type())
	}
}

func (i *Identifier) GenerateTypeEqual() {
	for nodeID, str := range i.NodeStrMap {
		originNodeID := i.StrNodeMap[str]
		if nodeID != originNodeID {
			dstStructName := i.NodeStructMap[nodeID]
			srcStructName := i.NodeStructMap[originNodeID]
			if dstStructName != srcStructName {
				i.StructEqualMap[dstStructName] = srcStructName
			}
		}
	}
	for dstStructName, srcStructName := range i.StructEqualMap {
		i.StructEquals = append(i.StructEquals, []string{dstStructName, srcStructName})
	}
	sort.Slice(i.StructEquals, func(m, n int) bool {
		return i.StructEquals[m][0] < i.StructEquals[n][0]
	})
}

func (i *Identifier) GenerateStr(node Node) {
	i.generateStr(node, func(node Node) {
		i.concat(node, false)
	})
}

func (i *Identifier) generateStr(node Node, f func(Node)) {
	for _, node := range node.Nodes() {
		i.generateStr(node, f)
	}
	if node.Type() == NodeTypeStruct {
		f(node)
	}
}

func (i *Identifier) concat(node Node, withName bool) (r string) {
	var s string
	switch node.Type() {
	case NodeTypeSimple:
		s = string(node.Field().Structure)
	case NodeTypeSlice:
		s = fmt.Sprintf("slice[%s]", i.concat(node.Nodes()[0], false))
	case NodeTypeMap:
		s = fmt.Sprintf("map[%s, %s]", i.concat(node.Nodes()[0], false), i.concat(node.Nodes()[1], false))
	case NodeTypeStruct:
		s = i.concatStruct(node)

		// add new type
		if _, ok := i.StrNodeMap[s]; !ok {
			i.StrNodeMap[s] = node.ID()
			i.OriginNodes = append(i.OriginNodes, node)
		}

		// record index
		i.NodeStrMap[node.ID()] = s
	default:
		Exit("[%v] Unsupported node type %s", node.Sheet(), node.Type())
	}
	if withName && node.Field() != nil && node.Field().Name != "" {
		s = fmt.Sprintf("%s(%s)", s, node.Field().Name)
	}

	return s
}

func (i *Identifier) concatStruct(node Node) string {
	var b strings.Builder
	var maxSize = len(node.Nodes()) - 1
	for index, node := range node.Nodes() {
		b.WriteString(i.concat(node, true))
		if index < maxSize {
			b.WriteString(FlagComma)
		}
	}
	return fmt.Sprintf("{%s}", b.String())
}
