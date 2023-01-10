package converter

import (
	"fmt"
)

type NodeMap struct {
	*NodeBase
}

func NewNodeMap(ctx *NodeContext, sheet Sheet, field *Field) *NodeMap {
	n := &NodeMap{
		NodeBase: NewNodeBase(ctx, sheet, field),
	}

	n.ParseNodes()
	return n
}

func (n *NodeMap) RootName() string {
	if n.ctx.key == "" {
		return n.NodeBase.RootName()
	}
	return fmt.Sprintf("%sKey%s", n.NodeBase.RootName(), n.ctx.key)
}

func (n *NodeMap) Type() NodeType {
	return NodeTypeMap
}

func (n *NodeMap) ParseNodes() {
	if n.field == nil {
		return
	}
	switch n.field.Structure {
	case StructureTypeTable:
		n.AppendNode(NewNode(n.ctx, n.sheet, n.field.Args[0]))
		n.AppendNode(NewNode(n.ctx, n.sheet, n.field.Args[1]))
	case StructureTypeMap:
		sheet := n.Excel().GetSheet(n.InferiorSheetName())
		sheet.ParseContent(n.field.Structure)
		if n.ctx.key == "" {
			n.AppendNode(NewNode(NewNodeContext(), sheet, 0))
		} else {
			headerField := n.sheet.GetHeaderField(n.ctx.key)
			n.AppendNode(NewNode(NewNodeContext(), sheet, headerField.Index))
		}
		n.AppendNode(NewNode(NewNodeContext(), sheet, StructureTypeStruct))
	case StructureTypeRowMap:
		sheet := n.Excel().GetSheet(n.InferiorSheetName())
		sheet.ParseContent(n.field.Structure)
		n.AppendNode(NewNode(NewNodeContext(), sheet, 0))
		n.AppendNode(NewNode(NewNodeContext(), sheet, StructureTypeRowMapVal))
	case StructureTypeColMap:
		sheet := n.Excel().GetSheet(n.InferiorSheetName())
		sheet.ParseContent(n.field.Structure)
		n.AppendNode(NewNode(NewNodeContext(), sheet, 0))
		n.AppendNode(NewNode(NewNodeContext(), sheet, StructureTypeColMapVal))
	default:
		Exit("[%v] Unsupported field type %s", n.sheet, n.field.Structure)
	}
}
