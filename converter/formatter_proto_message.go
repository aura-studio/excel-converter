package converter

import (
	"fmt"
)

type FormatterProtoMessages struct {
	*FormatterBase
	used       bool
	identifier *Identifier
}

func NewFormatterProtoMessages(identifier *Identifier) *FormatterProtoMessages {
	f := &FormatterProtoMessages{
		FormatterBase: NewFormatterBase(),
		identifier:    identifier,
	}

	// 写入proto文件头部
	f.WriteString(`syntax = "proto3";
`)
	return f
}

func (f *FormatterProtoMessages) FormatMessages() {
	f.used = true

	// 在消息定义开始前添加空行
	f.WriteString("\n")

	// 为每个结构生成消息定义
	for _, node := range f.identifier.OriginNodes {
		f.FormatMessage(node)
		f.WriteString("\n")
	}
}

func (f *FormatterProtoMessages) FormatMessage(node Node) {
	structName := f.identifier.NodeStructMap[node.ID()]

	f.WriteString(fmt.Sprintf("message %s {\n", structName))

	// 生成字段
	fieldNumber := 1
	for _, childNode := range node.Nodes() {
		f.WriteString("  ")
		f.WriteString(f.getProtoType(childNode))
		f.WriteString(" ")
		f.WriteString(f.getProtoFieldName(childNode.FieldName()))
		f.WriteString(" = ")
		f.WriteString(fmt.Sprintf("%d", fieldNumber))
		f.WriteString(";\n")
		fieldNumber++
	}

	f.WriteString("}\n")
}

// 保持Go风格的字段名（如ID而不是i_d）
func (f *FormatterProtoMessages) getProtoFieldName(fieldName string) string {
	// 直接返回原始字段名，保持Go命名风格
	return fieldName
}

// 根据节点类型获取对应的proto类型
func (f *FormatterProtoMessages) getProtoType(node Node) string {
	switch node.Type() {
	case NodeTypeSimple:
		return f.getProtoSimpleType(node)
	case NodeTypeSlice:
		return fmt.Sprintf("repeated %s", f.getProtoType(node.Nodes()[0]))
	case NodeTypeMap:
		keyType := f.getProtoType(node.Nodes()[0])
		valueType := f.getProtoType(node.Nodes()[1])
		return fmt.Sprintf("map<%s, %s>", keyType, valueType)
	case NodeTypeStruct:
		return f.identifier.NodeStructMap[node.ID()]
	default:
		Exit("[%v] Unsupported node type %s", node.Sheet(), node.Type())
		return ""
	}
}

// 获取基础类型的proto类型
func (f *FormatterProtoMessages) getProtoSimpleType(node Node) string {
	switch node.Field().Structure {
	case StructureTypeString:
		return "string"
	case StructureTypeInt:
		return "int64"
	case StructureTypeFloat:
		return "double"
	case StructureTypeBool:
		return "bool"
	case StructureTypeBigInt:
		return "string" // proto3 没有big int，使用string表示
	case StructureTypeBigFloat:
		return "string" // proto3 没有big float，使用string表示
	case StructureTypeBigRat:
		return "string" // proto3 没有big rat，使用string表示
	case StructureTypeTime:
		return "int64" // 使用timestamp (unix时间戳)
	default:
		Exit("[%v] Unsupported structure type %s", node.Sheet(), node.Field().Structure)
		return ""
	}
}

func (f *FormatterProtoMessages) Close() string {
	if !f.used {
		return ""
	}
	return f.String()
}
