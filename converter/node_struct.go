package converter

import (
	"github.com/spf13/cast"
)

type NodeStruct struct {
	*NodeBase
	verticleIndex int
}

func NewNodeStruct(ctx *NodeContext, sheet Sheet, field *Field) *NodeStruct {
	n := &NodeStruct{
		NodeBase: NewNodeBase(ctx, sheet, field),
	}
	n.ParseVerticleIndex()
	n.ParseNodes()
	return n
}

func (n *NodeStruct) Type() NodeType {
	return NodeTypeStruct
}

func (n *NodeStruct) ParseVerticleIndex() {
	if n.field == nil {
		return
	}
	if n.field.Structure == StructureTypeDict || n.field.Structure == StructureTypeList {
		if n.field.Index >= 0 {
			n.verticleIndex = n.field.Index + 1
		} else {
			n.verticleIndex = n.ctx.index
		}
	}
}

func (n *NodeStruct) ParseNodes() {
	if n.field == nil {
		return
	}
	verticleIndex := n.verticleIndex
	switch n.field.Structure {
	case StructureTypeDict:
		for _, typArg := range n.field.Args {
			if size, err := cast.ToIntE(typArg); err == nil {
				for index := 0; index < size; {
					field := NewField(n.sheet, n.sheet.GetHeaderField(verticleIndex))
					n.ctx.index = verticleIndex
					n.AppendNode(NewNode(n.ctx, n.sheet, verticleIndex))
					index++
					verticleIndex += field.Size
				}
			} else {
				field := NewField(n.sheet, typArg)
				n.ctx.index = verticleIndex
				n.AppendNode(NewNode(n.ctx, n.sheet, typArg))
				verticleIndex += field.Size
			}
		}
	case StructureTypeStruct, StructureTypeMapVal:
		sheet := n.Excel().GetSheet(n.InferiorSheetName())
		sheet.ParseContent(n.field.Structure)
		for index := 0; index < sheet.HeaderSize(); {
			field := NewField(sheet, sheet.GetHeaderField(verticleIndex))
			n.AppendNode(NewNode(NewNodeContext(), sheet, index))
			index += field.Size
			verticleIndex += field.Size
		}
	default:
		Exit("[%v] Unsupported field type %s", n.sheet, n.field.Structure)
	}
}
