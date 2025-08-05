package converter

import (
	"path/filepath"
	"sort"
)

type ConverterProto struct {
	*ConverterBase
	identifier *Identifier
	collection *Collection
}

func NewConverterProto() *ConverterProto {
	c := &ConverterProto{
		ConverterBase: NewConverterBase(ConverterTypeProto),
		identifier:    NewIdentifier(),
		collection:    NewCollection(),
	}
	return c
}

func (c *ConverterProto) Run() {
	c.Load()
	c.Parse()
	c.Export()
}

func (c *ConverterProto) Parse() {
	c.Build()
	c.Identity()
	c.Link()
}

func (c *ConverterProto) Export() {
	c.Format()
	c.Remove()
	c.Write()
}

func (c *ConverterProto) Identity() {
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

func (c *ConverterProto) Link() {
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

func (c *ConverterProto) Format() {
	c.FormatMessages()
}

func (c *ConverterProto) FormatMessages() {
	formatter := NewFormatterProtoMessages(c.identifier)
	formatter.FormatMessages()
	content := formatter.Close()
	filepath := filepath.Join(path.ExportAbsPath(), "messages.proto")
	c.contentMap[filepath] = content
}
