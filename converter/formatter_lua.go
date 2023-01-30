package converter

import (
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type FormatterLua struct {
	*FormatterBase
	used bool
}

func NewFormatterLua() *FormatterLua {
	f := &FormatterLua{
		FormatterBase: NewFormatterBase(),
	}
	f.WriteString(`-- <important: auto generate by excel-to-lua converter, do not modify>
local _ = {}

`)
	return f
}

func (f *FormatterLua) Close() string {
	if !f.used {
		return ""
	}
	f.WriteString(`return _`)
	return f.String()
}

func (f *FormatterLua) FormatNode(node Node) {
	f.used = true
	sheetName := node.InferiorSheetName()
	excel := node.Excel()
	sheet := excel.GetSheet(sheetName)
	ctx := NewSourceContext()
	ctx.key = node.Context().key
	source := NewSourceTable(ctx, sheet, node.Field().Structure, nil)
	f.WriteString("_.")
	f.formatVarName(node)
	f.WriteString(" = ")
	f.formatValue(node, source)
	f.WriteString("\n\n")
}

func (f *FormatterLua) formatVarName(node Node) {
	f.WriteString(node.RootName())
}

func (f *FormatterLua) formatFieldName(node Node) {
	f.WriteString(format.ToLower(node.FieldName()))
}

func (f *FormatterLua) formatValue(node Node, source Source) {
	if source.Type() == SourceTypeNil {
		return
	}

	switch node.Type() {
	case NodeTypeSimple:
		f.FormatBase(node, []Source{source})
	case NodeTypeSlice:
		f.FormatSlice(node, source.Sources())
	case NodeTypeMap:
		f.FormatMap(node, source.Sources())
	case NodeTypeStruct:
		f.FormatStruct(node, []Source{source})
	}
}

func (f *FormatterLua) FormatBase(node Node, sources []Source) {
	source := sources[0]
	switch node.Field().Structure {
	case StructureTypeString:
		f.WriteString("\"")
		f.WriteString(source.Content())
		f.WriteString("\"")
	case StructureTypeInt, StructureTypeFloat:
		switch source.Content() {
		case "":
			f.WriteString("0")
		default:
			f.WriteString(source.Content())
		}
	case StructureTypeBool:
		switch source.Content() {
		case "0", FlagFalse, "":
			f.WriteString(FlagFalse)
		case "1", FlagTrue:
			f.WriteString(FlagTrue)
		default:
			Exit("[%v] Unknown bool value %s", source.Sheet(), source.Content())
		}
	case StructureTypeBigInt:
		switch source.Content() {
		case "":
			f.WriteString("0")
		default:
			s, ok := format.BigIntToLua(source.Content())
			if !ok {
				Exit("[%v] Unknown bigint value %s", source.Sheet(), source.Content())
			}
			f.WriteString(s)
		}
	case StructureTypeBigFloat:
		switch source.Content() {
		case "":
			f.WriteString("0")
		default:
			s, ok := format.BigFloatToLua(source.Content())
			if !ok {
				Exit("[%v] Unknown bigfloat value %s", source.Sheet(), source.Content())
			}
			f.WriteString(s)
		}
	case StructureTypeBigRat:
		f.WriteString("nil")
	case StructureTypeTime:
		switch source.Content() {
		case "":
			f.WriteString("0")
		default:
			content := source.Content()
			var tm time.Time
			var err error
			if len(content) == 0 {
				tm = time.Time{}
			} else if !strings.Contains(content, ".") {
				tm, err = format.ParseTime(content)
				if err != nil {
					Exit("[%v] Unknown time format %s", source.Sheet(), content)
				}
			} else {
				cellFloat, err := strconv.ParseFloat(content, 64)
				if err != nil {
					Exit("[%v] Unknown time format %s", source.Sheet(), content)
				}
				tm, err = excelize.ExcelDateToTime(cellFloat, false)
				if err != nil {
					Exit("[%v] Unknown time format %s", source.Sheet(), content)
				}
			}
			f.WriteString(strconv.FormatInt(tm.Unix(), 10))
		}
	default:
		Exit("[%v] Unknown structure %s", source.Sheet(), node.Field().Structure)
	}
}

func (f *FormatterLua) FormatSlice(node Node, sources []Source) {
	nodeSub := node.Nodes()[0]
	if nodeSub.Type() == NodeTypeSimple {
		f.WriteString("{")
		for index, source := range sources {
			f.formatValue(nodeSub, source)
			if index < len(sources)-1 {
				f.WriteString(", ")
			}
		}
		f.WriteString("}")
	} else {
		f.WriteString("{\n")
		f.IncDepth()
		for _, source := range sources {
			if f.IsSourceEmpty(source) {
				continue
			}
			f.FormatIndent()
			f.formatValue(nodeSub, source)
			f.WriteString(",\n")
		}
		f.DecDepth()
		f.FormatIndent()
		f.WriteString("}")
	}
}

func (f *FormatterLua) FormatMap(node Node, sources []Source) {
	nodeKey := node.Nodes()[0]
	nodeVal := node.Nodes()[1]

	f.WriteString("{\n")
	f.IncDepth()
	for index := 0; index < len(sources); {
		f.FormatIndent()
		f.WriteString("[")
		f.formatValue(nodeKey, sources[index])
		f.WriteString("]")
		index++
		f.WriteString(" = ")
		f.formatValue(nodeVal, sources[index])
		f.WriteString(",\n")
		index++
	}
	f.DecDepth()
	f.FormatIndent()
	f.WriteString("}")
}

func (f *FormatterLua) FormatStruct(node Node, sources []Source) {
	nodes := node.Nodes()
	sources = sources[0].Sources()

	f.WriteString("{\n")
	f.IncDepth()
	for index, source := range sources {
		nodeSub := nodes[index]
		if source.Type() == SourceTypeNil && nodeSub.Type() != NodeTypeSimple {
			continue
		}
		f.FormatIndent()
		f.formatFieldName(nodeSub)
		f.WriteString(" = ")
		f.formatValue(nodeSub, source)
		f.WriteString(",\n")
	}
	f.DecDepth()
	f.FormatIndent()
	f.WriteString("}")
}
