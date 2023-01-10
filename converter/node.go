package converter

type NodeType string

const (
	NodeTypeSimple NodeType = "simple" // 基础类型
	NodeTypeSlice  NodeType = "slice"  // 顺序类型 Sequence StructNode
	NodeTypeMap    NodeType = "map"    // 键值类型 FloatingMapping StructNode
	NodeTypeStruct NodeType = "struct" // 结""构类型 FixedMapping []StructNode
)

type NodeContext struct {
	index int
	key   string
}

func NewNodeContext() *NodeContext {
	return &NodeContext{
		index: -1,
	}
}

type Node interface {
	ID() uint32
	FieldName() string
	RootName() string
	Type() NodeType
	String() string
	GetCell(int) string
	HeaderSize() int
	InferiorSheetName() string
	Excel() Excel
	Sheet() Sheet
	Nodes() []Node
	Field() *Field
	AppendNode(node Node)
	ParseNodes()
	ExcelPathName() string
	SheetPathName() string
	Context() *NodeContext
}

var nodeCreators map[NodeType]func(ctx *NodeContext, sheet Sheet, field *Field) Node

func init() {
	nodeCreators = map[NodeType]func(ctx *NodeContext, sheet Sheet, field *Field) Node{
		NodeTypeSimple: func(ctx *NodeContext, sheet Sheet, field *Field) Node {
			return NewNodeSimple(ctx, sheet, field)
		},
		NodeTypeMap: func(ctx *NodeContext, sheet Sheet, field *Field) Node {
			return NewNodeMap(ctx, sheet, field)
		},
		NodeTypeSlice: func(ctx *NodeContext, sheet Sheet, field *Field) Node {
			return NewNodeSlice(ctx, sheet, field)
		},
		NodeTypeStruct: func(ctx *NodeContext, sheet Sheet, field *Field) Node {
			return NewNodeStruct(ctx, sheet, field)
		},
	}
}

func NewNode(ctx *NodeContext, sheet Sheet, v interface{}) Node {
	var nodeType NodeType
	var field *Field
	switch v := v.(type) {
	case NodeType:
		nodeType = v
	case StructureType:
		field = NewField(sheet, v)
	case int:
		field = NewField(sheet, sheet.GetHeaderField(v))
	case string:
		field = NewField(sheet, v)
	}
	if nodeType == "" {
		switch field.Structure {
		case StructureTypeString, StructureTypeInt, StructureTypeFloat, StructureTypeBool,
			StructureTypeBigInt, StructureTypeBigRat, StructureTypeBigFloat, StructureTypeTime:
			nodeType = NodeTypeSimple
		case StructureTypeRow, StructureTypeCol, StructureTypeRowMapVal, StructureTypeColMapVal,
			StructureTypeArray, StructureTypeStructs, StructureTypeRows, StructureTypeCols, StructureTypeList:
			nodeType = NodeTypeSlice
		case StructureTypeTable, StructureTypeMap, StructureTypeRowMap, StructureTypeColMap:
			nodeType = NodeTypeMap
		case StructureTypeDict, StructureTypeStruct, StructureTypeMapVal:
			nodeType = NodeTypeStruct
		default:
			Exit("[%s] Unsupported field type [%s]", sheet, field.Structure)
		}
	}
	if nodeCreator, ok := nodeCreators[nodeType]; !ok {
		Exit("[%v] Unsupported node type %s", sheet, nodeType)
	} else {
		return nodeCreator(ctx, sheet, field)
	}
	return nil
}
