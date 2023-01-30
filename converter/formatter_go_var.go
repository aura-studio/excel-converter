package converter

import (
	"fmt"
)

type FormatterGoVar struct {
	*FormatterBase
	used        bool
	packageName string
	identifier  *Identifier
}

func NewFormatterGoVar(packageName string, identifier *Identifier) *FormatterGoVar {
	f := &FormatterGoVar{
		FormatterBase: NewFormatterBase(),
		packageName:   packageName,
		identifier:    identifier,
	}
	f.WriteString(fmt.Sprintf(`// <important: auto generate by excel-to-go converter, do not modify>
package %s

import "%s/structs"

var (
`, packageName, path.ImportPath()))
	f.IncDepth()
	return f
}

func (f *FormatterGoVar) Close() string {
	if !f.used {
		return ""
	}
	f.WriteString(")")
	return f.String()
}

func (f *FormatterGoVar) FormatNode(node Node) {
	f.used = true
	sheetName := node.InferiorSheetName()
	excel := node.Excel()
	sheet := excel.GetSheet(sheetName)
	ctx := NewSourceContext()
	ctx.key = node.Context().key
	source := NewSourceTable(ctx, sheet, node.Field().Structure, nil)
	f.WriteString("\t")
	f.FormatVarName(node)
	f.WriteString(" = ")
	f.FormatValue(node, source, nil)
	f.WriteString("\n")
}

func (f *FormatterGoVar) FormatVarName(node Node) {
	f.WriteString(node.RootName())
}

func (f *FormatterGoVar) FormatFieldName(node Node) {
	f.WriteString(format.ToUpper(node.FieldName()))
}

func (f *FormatterGoVar) FormatValue(node Node, source Source, parentNode Node) {
	if source.Type() == SourceTypeNil {
		return
	}

	switch node.Type() {
	case NodeTypeSimple:
		f.FormatBase(node, []Source{source})
	case NodeTypeSlice:
		if parentNode == nil || (parentNode.Type() != NodeTypeSlice && parentNode.Type() != NodeTypeMap) {
			f.WriteString(f.identifier.NodeDataTypeMap[node.ID()])
		}
		f.FormatSlice(node, source.Sources())
	case NodeTypeMap:
		f.WriteString(f.identifier.NodeDataTypeMap[node.ID()])
		f.FormatMap(node, source.Sources())
	case NodeTypeStruct:
		if parentNode == nil || parentNode.Type() == NodeTypeStruct {
			f.WriteString(fmt.Sprintf("&%s", f.identifier.NodeDataTypeMap[node.ID()][1:]))
		}
		f.FormatStruct(node, []Source{source})
	}
}

func (f *FormatterGoVar) FormatBase(node Node, sources []Source) {
	source := sources[0]
	switch node.Field().Structure {
	case StructureTypeString:
		f.WriteString(fmt.Sprintf(`"%s"`, source.Content()))
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
		f.WriteString("structs.NewBigInt(\"")
		f.WriteString(source.Content())
		f.WriteString("\")")
	case StructureTypeBigFloat:
		f.WriteString("structs.NewBigFloat(\"")
		f.WriteString(source.Content())
		f.WriteString("\")")
	case StructureTypeBigRat:
		f.WriteString("structs.NewBigRat(\"")
		f.WriteString(source.Content())
		f.WriteString("\")")
	default:
		f.WriteString("\"")
		f.WriteString(source.Content())
		f.WriteString("\"")
	}
}

func (f *FormatterGoVar) FormatSlice(node Node, sources []Source) {
	f.WriteString("{}")
}

func (f *FormatterGoVar) FormatMap(node Node, sources []Source) {
	f.WriteString("{}")
}

func (f *FormatterGoVar) FormatStruct(node Node, sources []Source) {
	f.WriteString("{}")
}
