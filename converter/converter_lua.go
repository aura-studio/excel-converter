package converter

import (
	"fmt"
	"path/filepath"
)

type ConverterLua struct {
	*ConverterBase
}

func NewConverterLua() *ConverterLua {
	c := &ConverterLua{
		ConverterBase: NewConverterBase(ConverterTypeLua),
	}
	return c
}

func (c *ConverterLua) Run() {
	c.Load()
	c.Parse()
	c.Export()
}

func (c *ConverterLua) Parse() {
	c.Build()
}

func (c *ConverterLua) Export() {
	c.Format()
	c.Remove()
	c.Write()
}

func (c *ConverterLua) Format() {
	domains := make([]Domain, 0)
	c.ForeachDomain(func(domain Domain) {
		domains = append(domains, domain)
	})
	results := c.Parallel(ToSlice(domains), func(param interface{}) func() interface{} {
		return func() interface{} {
			domain := param.(Domain)
			formatter := NewFormatterLua()
			for _, excel := range domain[ExcelTypeRegular] {
				for _, node := range excel.Nodes() {
					if node.Excel().ForClient() && node.Sheet().ForClient() {
						formatter.FormatNode(node)
					}
				}
			}
			content := formatter.Close()
			if len(content) == 0 {
				return nil
			}
			return []string{c.GetFilePath(domain), content}
		}
	})
	for _, result := range results {
		if result == nil {
			continue
		}
		filePath := result.([]string)[0]
		content := result.([]string)[1]
		c.contentMap[filePath] = content
	}
}

func (c *ConverterLua) GetFilePath(domain Domain) string {
	for _, excels := range domain {
		for _, excel := range excels {
			packageName := format.ToLuaPackageCase(excel.PackageName())
			fileName := fmt.Sprintf("%s.lua", excel.DomainName())
			return filepath.Join(path.ExportAbsPath(), packageName, fileName)
		}
	}
	Exit("[Main] Cannot find excel in domain")
	return ""
}
