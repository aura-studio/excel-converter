package converter

import "fmt"

type FormatterGoStorageTypes struct {
	*FormatterBase
	identifier *Identifier
}

func NewFormatterGoStorageTypes(identifier *Identifier) *FormatterGoStorageTypes {
	f := &FormatterGoStorageTypes{
		FormatterBase: NewFormatterBase(),
		identifier:    identifier,
	}
	f.WriteString(`// <important: auto generate by excel-to-go converter, do not modify>
package storage
`)
	return f
}

func (f *FormatterGoStorageTypes) FormatPackages() {
	f.WriteString("\nimport (\n")
	f.WriteString("\t\"")
	f.WriteString(path.ImportPath())
	f.WriteString("/structs\"\n")
	f.WriteString(")\n")
}

func (f *FormatterGoStorageTypes) FormatVars() {
	f.WriteString(`
var TypeStorage = make(map[string]map[string]any)
`)
}

func (f *FormatterGoStorageTypes) FormatFuncs() {
	f.WriteString(`
func LoadType(excelName, sheetName string, v any) {
	if _, ok := TypeStorage[excelName]; !ok {
		TypeStorage[excelName] = make(map[string]any)
	}
	TypeStorage[excelName][sheetName] = v
}

func LoadTypes() {
`)
}

func (f *FormatterGoStorageTypes) FormatNode(node Node) {
	sheetName := node.InferiorSheetName()
	excel := node.Excel()
	sheet := excel.GetSheet(sheetName)
	ctx := NewSourceContext()
	ctx.key = node.Context().key
	source := NewSourceTable(ctx, sheet, node.Field().Structure, nil)
	f.WriteString("\tLoadType(\"")
	f.WriteString(node.ExcelPathName())
	f.WriteString("\", \"")
	f.WriteString(node.SheetPathName())
	f.WriteString("\", ")
	f.FormatValue(node, source, nil)
	f.WriteString(")\n")
}

func (f *FormatterGoStorageTypes) FormatVarName(node Node) {
	f.WriteString(node.RootName())
}

func (f *FormatterGoStorageTypes) FormatFieldName(node Node) {
	f.WriteString(format.ToUpper(node.FieldName()))
}

func (f *FormatterGoStorageTypes) FormatValue(node Node, source Source, parentNode Node) {
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

func (f *FormatterGoStorageTypes) FormatBase(node Node, sources []Source) {
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

func (f *FormatterGoStorageTypes) FormatSlice(node Node, sources []Source) {
	f.WriteString("{}")
}

func (f *FormatterGoStorageTypes) FormatMap(node Node, sources []Source) {
	f.WriteString("{}")
}

func (f *FormatterGoStorageTypes) FormatStruct(node Node, sources []Source) {
	f.WriteString("{}")
}

func (f *FormatterGoStorageTypes) Close() string {
	f.WriteString("}")
	return f.String()
}
