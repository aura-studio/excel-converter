package converter

import (
	"fmt"
)

type FormatterGoVars struct {
	*FormatterBase
	used        bool
	packageName string
	identifier  *Identifier
}

func NewFormatterGoVar(packageName string, identifier *Identifier) *FormatterGoVars {
	f := &FormatterGoVars{
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

func (f *FormatterGoVars) FormatNode(node Node) {
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

func (f *FormatterGoVars) FormatVarName(node Node) {
	f.WriteString(node.RootName())
}

func (f *FormatterGoVars) FormatFieldName(node Node) {
	f.WriteString(format.ToUpper(node.FieldName()))
}

func (f *FormatterGoVars) FormatValue(node Node, source Source, parentNode Node) {
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

func (f *FormatterGoVars) FormatBase(node Node, sources []Source) {
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

func (f *FormatterGoVars) FormatSlice(node Node, sources []Source) {
	f.WriteString("{}")
}

func (f *FormatterGoVars) FormatMap(node Node, sources []Source) {
	f.WriteString("{}")
}

func (f *FormatterGoVars) FormatStruct(node Node, sources []Source) {
	f.WriteString("{}")
}

func (f *FormatterGoVars) Close() string {
	if !f.used {
		return ""
	}
	f.WriteString(")")
	return f.String()
}
