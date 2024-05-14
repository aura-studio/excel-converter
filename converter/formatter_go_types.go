package converter

import "fmt"

type FormatterGoTypes struct {
	*FormatterBase
	exelMap    map[string]map[string]bool
	identifier *Identifier
}

func NewFormatterGoTypes(identifier *Identifier) *FormatterGoTypes {
	f := &FormatterGoTypes{
		FormatterBase: NewFormatterBase(),
		exelMap:       make(map[string]map[string]bool),
	}
	f.WriteString(`// <important: auto generate by excel-to-go converter, do not modify>
package storage

`)
	return f
}

func (f *FormatterGoTypes) FormatVars() {
	f.WriteString(`
var TypeMap = make(map[string]map[string]any)
`)
}

func (f *FormatterGoTypes) FormatFuncs() {
	f.WriteString(`
func Register(excelName, sheetName string, v any) {
	if _, ok := TypeStorage[excelName]; !ok {
		TypeMap[excelName] = make(map[string]any)
	}
	TypeMap[excelName][sheetName] = v
}

func CategoriesLoading() {
`)
}

func (f *FormatterGoTypes) FormatCategories(categories []string) {
	for _, category := range categories {
		f.WriteString("\tCategory(\"")
		f.WriteString(category)
		f.WriteString("\"),\n")
	}
	f.WriteString("}\n")
}

func (f *FormatterGoTypes) FormatNode(node Node) {
	sheetName := node.InferiorSheetName()
	excel := node.Excel()
	sheet := excel.GetSheet(sheetName)
	ctx := NewSourceContext()
	ctx.key = node.Context().key
	source := NewSourceTable(ctx, sheet, node.Field().Structure, nil)
	f.WriteString("\t")
	f.WriteString("\tRegister(\"")
	f.WriteString(excel.Name())
	f.WriteString("\", \"")
	f.WriteString(sheetName)
	f.WriteString("\", ")
	f.FormatValue(node, source, nil)
	f.WriteString("\")\n")
	f.WriteString("\n")
}

func (f *FormatterGoTypes) FormatVarName(node Node) {
	f.WriteString(node.RootName())
}

func (f *FormatterGoTypes) FormatFieldName(node Node) {
	f.WriteString(format.ToUpper(node.FieldName()))
}

func (f *FormatterGoTypes) FormatValue(node Node, source Source, parentNode Node) {
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

func (f *FormatterGoTypes) FormatBase(node Node, sources []Source) {
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

func (f *FormatterGoTypes) FormatSlice(node Node, sources []Source) {
	f.WriteString("{}")
}

func (f *FormatterGoTypes) FormatMap(node Node, sources []Source) {
	f.WriteString("{}")
}

func (f *FormatterGoTypes) FormatStruct(node Node, sources []Source) {
	f.WriteString("{}")
}

func (f *FormatterGoTypes) Close() string {
	f.WriteString(")")
	return f.String()
}
