package converter

import (
	"fmt"
	"path/filepath"
)

type RenderJson struct{}

func NewRenderJson() *RenderJson {
	return &RenderJson{}
}

func (r *RenderJson) Render() {
	domains := make([]Domain, 0)
	c.ForeachDomain(func(domain Domain) {
		domains = append(domains, domain)
	})
	results := c.Parallel(ToSlice(domains), func(param any) func() any {
		return func() any {
			domain := param.(Domain)
			formatter := NewFormatterJSON(r.GetPackageName(domain))
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

func (r *RenderJson) GetFilePath(domain Domain) string {
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

func (r *RenderJson) GetPackageName(domain Domain) string {
	for _, excels := range domain {
		for _, excel := range excels {
			return format.ToJsonPackageCase(excel.PackageName())
		}
	}
	Exit("[Main] Cannot find excel in domain")
	return ""
}
