package converter

import (
	"fmt"
	"strconv"
)

type FormatterGoJSON struct {
	*FormatterBase
	used        bool
	packageName string
	identifier  *Identifier
}

func NewFormatterGoJSON(packageName string, identifier *Identifier) *FormatterGoJSON {
	f := &FormatterGoJSON{
		FormatterBase: NewFormatterBase(),
		packageName:   packageName,
		identifier:    identifier,
	}

	f.WriteString(fmt.Sprintf(`// <important: auto generate by excel-to-go converter, do not modify>
package %s

import "encoding/json"

func init() {
	var (
		data string
		err error
	)

`, packageName))
	f.IncDepth()
	return f
}

func (f *FormatterGoJSON) Close() string {
	if !f.used {
		return ""
	}
	f.WriteString("}")
	return f.String()
}

func (f *FormatterGoJSON) FormatNode(node Node) {
	f.used = true
	sheetName := node.InferiorSheetName()
	excel := node.Excel()
	sheet := excel.GetSheet(sheetName)
	ctx := NewSourceContext()
	ctx.key = node.Context().key
	source := NewSourceTable(ctx, sheet, node.Field().Structure, nil)
	f.WriteString("\t")
	f.WriteString("data = `")
	f.FormatValue(node, source, nil)
	f.WriteString("`\n\n")
	f.WriteString("\tif err = json.Unmarshal([]byte(data), &")
	f.FormatVarName(node)
	f.WriteString(`); err != nil {
		panic(err)
	}`)
	f.WriteString("\n\n")
}

func (f *FormatterGoJSON) FormatVarName(node Node) {
	f.WriteString(node.RootName())
}

func (f *FormatterGoJSON) FormatFieldName(node Node) {
	f.WriteString("\"")
	f.WriteString(format.ToUpper(node.FieldName()))
	f.WriteString("\"")
}

func (f *FormatterGoJSON) FormatValue(node Node, source Source, parentNode Node) {
	switch node.Type() {
	case NodeTypeSimple:
		f.FormatBase(node, []Source{source})
	case NodeTypeSlice:
		f.FormatSlice(node, source.Sources())
	case NodeTypeMap:
		f.FormatMap(node, source.Sources())
	case NodeTypeStruct:
		f.FormatStruct(node, []Source{source})
	default:
		Exit("[%v] Unknown structure %s", source.Sheet(), node.Field().Structure)
	}
}

func (f *FormatterGoJSON) FormatBase(node Node, sources []Source) {
	source := sources[0]
	switch node.Field().Structure {
	case StructureTypeString:
		f.WriteString(strconv.Quote(source.Content()))
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
	case StructureTypeBigInt: // TODO: not support yet. TO Fill func MarshalJSON UnmarshalJSON
		f.WriteString(source.Content())
	case StructureTypeBigFloat: // TODO: not support yet. TO Fill func MarshalJSON UnmarshalJSON
		f.WriteString(source.Content())
	case StructureTypeBigRat: // TODO: not support yet. TO Fill func MarshalJSON UnmarshalJSON
		f.WriteString(source.Content())
	case StructureTypeTime: // TODO: not support yet. TO Fill func MarshalJSON UnmarshalJSON
		f.WriteString(source.Content())
	default:
		Exit("[%v] Unknown structure %s", source.Sheet(), node.Field().Structure)
	}
}

func (f *FormatterGoJSON) FormatFieldBase(node Node, sources []Source) {
	f.WriteByte('"')
	f.FormatBase(node, sources)
	f.WriteByte('"')
}

func (f *FormatterGoJSON) FormatSlice(node Node, sources []Source) {
	nodeSub := node.Nodes()[0]
	if nodeSub.Type() == NodeTypeSimple {
		f.WriteString("[")
		for index, source := range sources {
			f.FormatValue(nodeSub, source, node)
			if index < len(sources)-1 {
				f.WriteString(",")
			}
		}
		f.WriteString("]")
	} else {
		f.WriteString("[\n")
		f.IncDepth()
		var validIndexes = make([]int, 0, len(sources))
		for index, source := range sources {
			if f.IsSourceEmpty(source) {
				continue
			}
			validIndexes = append(validIndexes, index)
		}
		for _, index := range validIndexes {
			source := sources[index]
			f.FormatIndent()
			f.FormatValue(nodeSub, source, node)
			if index < len(validIndexes)-1 {
				f.WriteString(",\n")
			} else {
				f.WriteString("\n")
			}
		}
		f.DecDepth()
		f.FormatIndent()
		f.WriteString("]")
	}
}

func (f *FormatterGoJSON) FormatMap(node Node, sources []Source) {
	nodes := node.Nodes()
	nodeKey := nodes[0]
	nodeVal := nodes[1]
	f.WriteString("{\n")
	f.IncDepth()
	for index := 0; index < len(sources); {
		f.FormatIndent()
		f.FormatFieldBase(nodeKey, []Source{sources[index]})
		index++
		f.WriteString(": ")
		f.FormatValue(nodeVal, sources[index], node)
		index++
		if index < len(sources)-1 {
			f.WriteString(",\n")
		} else {
			f.WriteString("\n")
		}
	}
	f.DecDepth()
	f.FormatIndent()
	f.WriteString("}")
}

func (f *FormatterGoJSON) FormatStruct(node Node, sources []Source) {
	nodes := node.Nodes()
	sources = sources[0].Sources()
	f.WriteString("{\n")
	f.IncDepth()
	var validIndexes = make([]int, 0, len(sources))
	for index, source := range sources {
		nodeSub := nodes[index]
		if source.Type() == SourceTypeNil && nodeSub.Type() != NodeTypeSimple {
			continue
		}
		validIndexes = append(validIndexes, index)
	}
	for _, index := range validIndexes {
		source := sources[index]
		nodeSub := nodes[index]
		if source.Type() == SourceTypeNil && nodeSub.Type() != NodeTypeSimple {
			continue
		}
		f.FormatIndent()
		f.FormatFieldName(nodeSub)
		f.WriteString(": ")
		f.FormatValue(nodeSub, source, node)
		if index < len(validIndexes)-1 {
			f.WriteString(",\n")
		} else {
			f.WriteString("\n")
		}
	}
	f.DecDepth()
	f.FormatIndent()
	f.WriteString("}")
}
