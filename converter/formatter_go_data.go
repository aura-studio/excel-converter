package converter

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type FormatterGoData struct {
	*FormatterBase
	used        bool
	packageName string
	identifier  *Identifier
}

func NewFormatterGoData(packageName string, identifier *Identifier) *FormatterGoData {
	f := &FormatterGoData{
		FormatterBase: NewFormatterBase(),
		packageName:   packageName,
		identifier:    identifier,
	}

	f.WriteString(fmt.Sprintf(`//go:build !debug
// +build !debug

// <important: auto generate by excel-to-go converter, do not modify>
package %s

import "%s/structs"

func init() {
`, packageName, path.ImportPath()))
	f.IncDepth()
	return f
}

func (f *FormatterGoData) Close() string {
	if !f.used {
		return ""
	}
	f.WriteString("}")
	return f.String()
}

func (f *FormatterGoData) FormatNode(node Node) {
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
	f.WriteString("\n\n")
}

func (f *FormatterGoData) FormatVarName(node Node) {
	f.WriteString(node.RootName())
}

func (f *FormatterGoData) FormatFieldName(node Node) {
	f.WriteString(format.ToUpper(node.FieldName()))
}

func (f *FormatterGoData) FormatValue(node Node, source Source, parentNode Node) {
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
	default:
		Exit("[%v] Unknown structure %s", source.Sheet(), node.Field().Structure)
	}
}

func (f *FormatterGoData) FormatBase(node Node, sources []Source) {
	source := sources[0]
	switch node.Field().Structure {
	case StructureTypeString:
		if strings.Contains(source.Content(), "\\") {
			f.WriteString(fmt.Sprintf("`%s`", source.Content()))
		} else {
			f.WriteString(fmt.Sprintf(`"%s"`, source.Content()))
		}
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
	case StructureTypeTime:
		f.WriteString("structs.NewTime(")
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
		f.WriteString(strconv.FormatInt(int64(tm.Year()), 10))
		f.WriteString(", ")
		f.WriteString(strconv.FormatInt(int64(tm.Month()), 10))
		f.WriteString(", ")
		f.WriteString(strconv.FormatInt(int64(tm.Day()), 10))
		f.WriteString(", ")
		f.WriteString(strconv.FormatInt(int64(tm.Hour()), 10))
		f.WriteString(", ")
		f.WriteString(strconv.FormatInt(int64(tm.Minute()), 10))
		f.WriteString(", ")
		f.WriteString(strconv.FormatInt(int64(tm.Second()), 10))
		f.WriteString(")")
	default:
		Exit("[%v] Unknown structure %s", source.Sheet(), node.Field().Structure)
	}
}

func (f *FormatterGoData) FormatSlice(node Node, sources []Source) {
	nodeSub := node.Nodes()[0]
	if nodeSub.Type() == NodeTypeSimple {
		f.WriteString("{")
		for index, source := range sources {
			f.FormatValue(nodeSub, source, node)
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
			f.FormatValue(nodeSub, source, node)
			f.WriteString(",\n")
		}
		f.DecDepth()
		f.FormatIndent()
		f.WriteString("}")
	}
}

func (f *FormatterGoData) FormatMap(node Node, sources []Source) {
	nodes := node.Nodes()
	nodeKey := nodes[0]
	nodeVal := nodes[1]
	f.WriteString("{\n")
	f.IncDepth()
	for index := 0; index < len(sources); {
		f.FormatIndent()
		f.FormatValue(nodeKey, sources[index], node)
		index++
		f.WriteString(": ")
		f.FormatValue(nodeVal, sources[index], node)
		f.WriteString(",\n")
		index++
	}
	f.DecDepth()
	f.FormatIndent()
	f.WriteString("}")
}

func (f *FormatterGoData) FormatStruct(node Node, sources []Source) {
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
		f.FormatFieldName(nodeSub)
		f.WriteString(": ")
		f.FormatValue(nodeSub, source, node)
		f.WriteString(",\n")
	}
	f.DecDepth()
	f.FormatIndent()
	f.WriteString("}")
}
