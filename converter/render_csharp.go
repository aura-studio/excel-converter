package converter

import (
	"fmt"
	"path/filepath"
)

type RenderCSharp struct{}

func NewRenderCSharp() *RenderCSharp {
	return &RenderCSharp{}
}

func (r *RenderCSharp) Render() {
	r.FormatStructs()
	r.FormatJson()
}

func (r *RenderCSharp) FormatStructs() {
	formatter := NewFormatterCSharpStructs(c.identifier)
	formatter.FormatStruct()
	formatter.FormatStructEqual()
	content := formatter.Close()
	filePath := filepath.Join(path.ExportAbsPath(), "Structs.cs")
	c.contentMap[filePath] = content
}

func (r *RenderCSharp) FormatJson() {
	domains := make([]Domain, 0)
	c.ForeachDomain(func(domain Domain) {
		domains = append(domains, domain)
	})
	results := c.Parallel(ToSlice(domains), func(param any) func() any {
		return func() any {
			domain := param.(Domain)
			formatter := NewFormatterCSharpJSON(r.GetPackageName(domain))
			for excelIdx, excel := range domain[ExcelTypeRegular] {
				nodes := make([]Node, 0)
				for _, node := range excel.Nodes() {
					if c.FilterNodeByDataType(node) {
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
			return []string{r.GetFilePath(domain), content}
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

func (r *RenderCSharp) GetFilePath(domain Domain) string {
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

func (r *RenderCSharp) GetPackageName(domain Domain) string {
	for _, excels := range domain {
		for _, excel := range excels {
			return format.ToJsonPackageCase(excel.PackageName())
		}
	}
	Exit("[Main] Cannot find excel in domain")
	return ""
}
