package converter

import (
	"github.com/spf13/cast"
)

type NodeSlice struct {
	*NodeBase
	verticleIndex int
}

func NewNodeSlice(ctx *NodeContext, sheet Sheet, field *Field) *NodeSlice {
	n := &NodeSlice{
		NodeBase: NewNodeBase(ctx, sheet, field),
	}
	n.ParseVerticleIndex()
	n.ParseNodes()
	return n
}

func (n *NodeSlice) Type() NodeType {
	return NodeTypeSlice
}

func (n *NodeSlice) ParseVerticleIndex() {
	if n.field == nil {
		return
	}
	if n.field.Structure == StructureTypeDict || n.field.Structure == StructureTypeList {
		if n.field.Index >= 0 {
			n.verticleIndex = n.field.Index + 1
		} else {
			n.verticleIndex = n.ctx.index
		}
		if n.field.Name == "" {
			n.field.Name = n.Sheet().GetHeaderField(n.verticleIndex).Name
		}
	}
}

func (n *NodeSlice) ParseNodes() {
	if n.field == nil {
		return
	}
	switch n.field.Structure {
	case StructureTypeArray:
		n.AppendNode(NewNode(n.ctx, n.sheet, n.field.Args[0]))
	case StructureTypeRow, StructureTypeCol:
		sheet := n.Excel().GetSheet(n.InferiorSheetName())
		sheet.ParseContent(n.field.Structure)
		n.AppendNode(NewNode(NewNodeContext(), sheet, 0))
	case StructureTypeRowMapVal, StructureTypeColMapVal:
		sheet := n.Excel().GetSheet(n.InferiorSheetName())
		sheet.ParseContent(n.field.Structure)
		n.AppendNode(NewNode(NewNodeContext(), sheet, 1))
	case StructureTypeRows, StructureTypeCols:
		sheet := n.Excel().GetSheet(n.InferiorSheetName())
		sheet.ParseContent(n.field.Structure)
		n.AppendNode(NewNode(NewNodeContext(), sheet, StructureTypeRow))
	case StructureTypeStructs:
		sheet := n.Excel().GetSheet(n.InferiorSheetName())
		sheet.ParseContent(n.field.Structure)
		n.AppendNode(NewNode(NewNodeContext(), sheet, StructureTypeStruct))
	case StructureTypeList:
		verticleIndex := n.verticleIndex
		field := NewField(n.Sheet(), n.field.Args[0])
		for index := 0; index < cast.ToInt(n.field.Args[1]); {
			switch field.Structure {
			case StructureTypeList, StructureTypeDict:
				n.ctx.index = verticleIndex
				n.AppendNode(NewNode(n.ctx, n.sheet, n.field.Args[0]))
			default:
				n.AppendNode(NewNode(n.ctx, n.sheet, verticleIndex))
			}
			index++
			verticleIndex += field.Size
		}

	default:
		Exit("[%v] Unsupported field type %s", n.sheet, n.field.Structure)
	}
}
