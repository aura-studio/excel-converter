package converter

import (
	"fmt"
	"strings"
	"sync/atomic"
)

var nodeIncrease uint32

type NodeBase struct {
	id    uint32
	field *Field
	ctx   *NodeContext
	sheet Sheet
	nodes []Node
}

func NewNodeBase(ctx *NodeContext, sheet Sheet, field *Field) *NodeBase {
	n := &NodeBase{
		id:    atomic.AddUint32(&nodeIncrease, 1),
		field: field,
		ctx:   ctx,
		sheet: sheet,
	}
	return n
}

func (n *NodeBase) ID() uint32 {
	return n.id
}

func (n *NodeBase) String() string {
	return fmt.Sprintf("%s:%s:%s:%s:%s", n.sheet.Excel().PackageName(), n.sheet.Excel().DomainName(),
		n.sheet.Excel().IndirectName(), n.sheet.IndirectName(), n.ctx.key)
}

func (n *NodeBase) FieldName() string {
	if n == nil {
		return ""
	}
	return n.field.Name
}

func (n *NodeBase) RootName() string {
	rootName := fmt.Sprintf("%s%s%s", format.ToUpper(n.Excel().DomainName()), format.ToUpper(n.Excel().IndirectName()), format.ToUpper(n.Sheet().IndirectName()))
	if rootName == "" {
		return FlagDefault
	}
	return rootName
}

func (n *NodeBase) Type() string {
	return ""
}

func (n *NodeBase) GetCell(index int) string {
	if n.field == nil {
		Exit("[%v] Structure %s can not get cell", n.sheet, n.field.Structure)
	}
	if n.field.Name != "" {
		return n.sheet.GetCell(n.field.HeaderField, index)
	}
	switch n.field.Structure {
	case StructureTypeArray:
	case StructureTypeTable:
	default:
		Exit("[%v] Structure %s can not get cell", n.sheet, n.field.Structure)
	}
	return ""
}

func (n *NodeBase) HeaderSize() int {
	if n.field == nil {
		Exit("[%v] Structure %s can not get cell", n.sheet, n.field.Structure)
	}
	if n.field.Name != "" {
		return n.sheet.HeaderSize()
	}
	return 0
}

func (n *NodeBase) InferiorSheetName() string {
	if n.field.HeaderField.Name == "" {
		if n.ctx.index < 0 { // Case 1: root node has node field
			return n.sheet.Name()
		}
		switch n.field.Structure {
		case StructureTypeRow, StructureTypeRows, StructureTypeCol, StructureTypeCols,
			StructureTypeMapVal, StructureTypeMap, StructureTypeStruct, StructureTypeStructs,
			StructureTypeRowMapVal, StructureTypeRowMap, StructureTypeColMapVal, StructureTypeColMap:
			content := n.sheet.GetCell(n.sheet.GetHeaderField(n.ctx.index), 0)
			strs := strings.Split(content, FlagComma)
			if len(strs) == 1 { // Case 2: normal inferior sheet
				return strs[0]
			}
			// Case 3: inferior sheet in sub nodes of list or dict
			strs = strings.Split(strs[0], FlagColon)
			if len(strs) == 1 {
				return strs[0]
			} else {
				return strs[1]
			}
		default: // Case 4: other type
			return n.sheet.Name()
		}
	} else if len(n.field.Args) > 0 {
		return n.field.Args[0]
	} else {
		return n.sheet.GetCell(n.field.HeaderField, 0)
	}
}

func (n *NodeBase) Excel() Excel {
	return n.sheet.Excel()
}

func (n *NodeBase) Sheet() Sheet {
	return n.sheet
}

func (n *NodeBase) Nodes() []Node {
	return n.nodes
}

func (n *NodeBase) Field() *Field {
	return n.field
}

func (n *NodeBase) AppendNode(node Node) {
	n.nodes = append(n.nodes, node)
}

func (n *NodeBase) ExcelPathName() string {
	if n.Excel().IndirectName() == FlagDefault {
		return n.Excel().DomainName()
	}
	return fmt.Sprintf("%s/%s", n.Excel().DomainName(), n.Excel().IndirectName())
}

func (n *NodeBase) SheetPathName() string {
	if n.ctx.key == "" {
		return n.Sheet().IndirectName()
	}
	return fmt.Sprintf("%s/%s", n.Sheet().IndirectName(), n.ctx.key)
}

func (n *NodeBase) Context() *NodeContext {
	return n.ctx
}
