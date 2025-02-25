package converter

import (
	"fmt"
	"path/filepath"
)

type ConverterJson struct {
	*ConverterBase
}

func NewConverterJson() *ConverterJson {
	c := &ConverterJson{
		ConverterBase: NewConverterBase(ConverterTypeJson),
	}
	return c
}

func (c *ConverterJson) Run() {
	c.Load()
	c.Parse()
	c.Export()
}

func (c *ConverterJson) Parse() {
	c.Build()
}

func (c *ConverterJson) Export() {
	c.Format()
	c.Remove()
	c.Write()
}

func (c *ConverterJson) GetJSONDataPath(domain Domain) string {
	for _, excels := range domain {
		for _, excel := range excels {
			goPackageName := format.ToGoPackageCase(excel.PackageName())
			fileName := fmt.Sprintf("%s.json", format.ToLower(excel.DomainName()))
			return filepath.Join(path.ExportAbsPath(), goPackageName, fileName)
		}
	}
	Exit("[Main] Cannot find excel in domain")
	return ""
}
func (c *ConverterJson) Format() {
	domains := make([]Domain, 0)
	c.ForeachDomain(func(domain Domain) {
		domains = append(domains, domain)
	})
	results := c.Parallel(ToSlice(domains), func(param any) func() any {
		return func() any {
			domain := param.(Domain)
			formatter := NewFormatterJSON(c.GetPackageName(domain))
			for excelIdx, excel := range domain[ExcelTypeRegular] {
				nodes := make([]Node, 0)
				for _, node := range excel.Nodes() {
					if node.Excel().ForClient() && node.Sheet().ForClient() {
						nodes = append(nodes, node)
					}
				}
				for i, node := range nodes {
					formatter.FormatNode(node, excelIdx == len(domain[ExcelTypeRegular])-1 && i == len(nodes)-1)
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

func (c *ConverterJson) GetFilePath(domain Domain) string {
	for _, excels := range domain {
		for _, excel := range excels {
			packageName := format.ToJsonPackageCase(excel.PackageName())
			fileName := fmt.Sprintf("%s.json", excel.DomainName())
			return filepath.Join(path.ExportAbsPath(), packageName, fileName)
		}
	}
	Exit("[Main] Cannot find excel in domain")
	return ""
}

func (c *ConverterJson) GetPackageName(domain Domain) string {
	for _, excels := range domain {
		for _, excel := range excels {
			return format.ToJsonPackageCase(excel.PackageName())
		}
	}
	Exit("[Main] Cannot find excel in domain")
	return ""
}
