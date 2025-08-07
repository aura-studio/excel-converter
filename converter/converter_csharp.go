package converter

import (
	"path/filepath"
	"sort"
)

type ConverterCSharp struct {
	*ConverterBase
	identifier *Identifier
	collection *Collection
}

func NewConverterCSharp() *ConverterCSharp {
	c := &ConverterCSharp{
		ConverterBase: NewConverterBase(ConverterTypeCSharp, FieldTypeClient),
		identifier:    NewIdentifier(),
		collection:    NewCollection(),
	}
	return c
}

func (c *ConverterCSharp) Run() {
	c.Load()
	c.Parse()
	c.Export()
}

func (c *ConverterCSharp) Parse() {
	c.Build()
	c.Identity()
	c.Link()
}

func (c *ConverterCSharp) Export() {
	c.Format()
	c.Remove()
	c.Write()
}

func (c *ConverterCSharp) Identity() {
	nodes := []Node{}
	c.ForeachExcel(func(e Excel) {
		if e.Type() == ExcelTypeTemplate {
			for _, node := range e.Nodes() {
				if node.Excel().ForClient() && node.Sheet().ForClient() {
					nodes = append(nodes, node)
				}
			}
		}
	})
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].String() < nodes[j].String()
	})
	for _, node := range nodes {
		c.identifier.GenerateStr(node)
	}
	nodes = []Node{}
	c.ForeachExcel(func(e Excel) {
		if e.Type() == ExcelTypeRegular {
			for _, node := range e.Nodes() {
				if node.Excel().ForClient() && node.Sheet().ForClient() {
					nodes = append(nodes, node)
				}
			}
		}
	})
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].String() < nodes[j].String()
	})
	for _, node := range nodes {
		c.identifier.GenerateStr(node)
	}
	c.ForeachExcel(func(e Excel) {
		if e.Type() == ExcelTypeTemplate || e.Type() == ExcelTypeRegular {
			for _, node := range e.Nodes() {
				if node.Excel().ForClient() && node.Sheet().ForClient() {
					c.identifier.GenerateStruct(node)
				}
			}
		}
	})
	c.ForeachExcel(func(e Excel) {
		if e.Type() == ExcelTypeTemplate || e.Type() == ExcelTypeRegular {
			for _, node := range e.Nodes() {
				if node.Excel().ForClient() && node.Sheet().ForClient() {
					c.identifier.GenerateType(node)
				}
			}
		}
	})

	c.identifier.GenerateTypeEqual()

	for str, nodeID := range c.identifier.StrNodeMap {
		Debug("[Identifier] struct[%v] = %s\n", nodeID, str)
	}
}

func (c *ConverterCSharp) Link() {
	c.ForeachExcel(func(e Excel) {
		if e.Type() == ExcelTypeRegular {
			for _, node := range e.Nodes() {
				if node.Excel().ForClient() && node.Sheet().ForClient() {
					c.collection.ReadNode(node)
				}
			}
		}
	})
	c.ForeachExcel(func(e Excel) {
		if e.Type() == ExcelTypeSettings {
			for _, sheets := range e.SheetMap() {
				for _, sheet := range sheets {
					c.collection.ReadLink(sheet)
				}
			}
		}
	})
}

func (c *ConverterCSharp) Format() {
	c.FormatStructs()
}

func (c *ConverterCSharp) FormatStructs() {
	formatter := NewFormatterCSharpStructs(c.identifier)
	formatter.FormatStruct()
	formatter.FormatStructEqual()
	content := formatter.Close()
	filePath := filepath.Join(path.ExportAbsPath(), "Structs.cs")
	c.contentMap[filePath] = content
}
