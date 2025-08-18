package converter

import (
	"strings"
)

type FormatterCSharpStructs struct {
	*FormatterBase
	used       bool
	identifier *Identifier
}

func NewFormatterCSharpStructs(identifier *Identifier) *FormatterCSharpStructs {
	f := &FormatterCSharpStructs{
		FormatterBase: NewFormatterBase(),
		identifier:    identifier,
	}
	f.WriteString(`// <important: auto generate by excel-to-csharp converter, do not modify>
using System;
using System.Collections.Generic;

namespace Exported.Config
{
`)

	return f
}

func (f *FormatterCSharpStructs) FormatStruct() {
	f.used = true
	for _, node := range f.identifier.OriginNodes {
		f.WriteString("    [Serializable]\n")
		f.WriteString("    public class ")
		// Add prefix to avoid conflict with field names
		f.WriteString("C_")
		f.WriteString(f.identifier.NodeStructMap[node.ID()])
		f.WriteString("\n    {\n")
		for _, childNode := range node.Nodes() {
			f.WriteString("        public ")
			goType := f.identifier.NodeTypeMap[childNode.ID()]
			csharpType := f.Translate(goType)
			f.WriteString(csharpType)
			f.WriteString(" ")
			f.WriteString(childNode.FieldName())
			f.WriteString(";\n")
		}
		f.WriteString("    }\n\n")
	}
}

func (f *FormatterCSharpStructs) FormatStructEqual() {
	f.used = true
	// C# doesn't have type aliases like Go, so we'll create inheritance relationships
	for _, structNames := range f.identifier.StructEquals {
		f.WriteString("    // Note: ")
		f.WriteString(structNames[0])
		f.WriteString(" is equivalent to ")
		f.WriteString(structNames[1])
		f.WriteString("\n")
		f.WriteString("    [Serializable]\n")
		f.WriteString("    public class ")
		f.WriteString("C_")
		f.WriteString(structNames[0])
		f.WriteString(" : ")
		f.WriteString("C_")
		f.WriteString(structNames[1])
		f.WriteString("\n    {\n    }\n\n")
	}
}

func (f *FormatterCSharpStructs) Close() string {
	if !f.used {
		return ""
	}
	f.WriteString("}")
	return f.String()
}

func (f *FormatterCSharpStructs) Translate(goType string) string {
	// 处理指针类型
	if strings.HasPrefix(goType, "*") {
		baseType := goType[1:]
		return f.Translate(baseType)
	}

	// 处理切片类型
	if strings.HasPrefix(goType, "[]") {
		elementType := goType[2:]
		return "List<" + f.Translate(elementType) + ">"
	}

	// 处理映射类型
	if strings.HasPrefix(goType, "map[") {
		// 解析 map[keyType]valueType
		mapContent := goType[4:] // 去掉 "map[" 和最后的 "]"
		// 找到键和值类型的分隔点
		splitIndex := -1
		for i, char := range mapContent {
			if char == ']' {
				splitIndex = i
				break
			}
		}
		if splitIndex > 0 {
			keyType := mapContent[:splitIndex]
			valueType := mapContent[splitIndex+2:]
			if valueType[0] == '*' {
				valueType = valueType[1:]
			}
			return "Dictionary<" + f.Translate(keyType) + ", " + f.Translate(valueType) + ">"
		}
	}

	// 基础类型映射
	switch goType {
	case "string":
		return "string"
	case "int64", "int", "int32":
		return "int"
	case "float64", "float32":
		return "float"
	case "bool":
		return "bool"
	case "big.Int":
		return "long"
	case "big.Float":
		return "decimal"
	case "big.Rat":
		return "decimal"
	case "time.Time":
		return "DateTime"
	default:
		// 如果是自定义结构体类型，直接返回
		return goType
	}
}
