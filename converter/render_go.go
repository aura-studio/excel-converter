package converter

import (
	"fmt"
	"path/filepath"
	"sort"
)

type RenderGo struct{}

func NewRenderGo() *RenderGo {
	return &RenderGo{}
}

func (r *RenderGo) Render() {
	r.FormatVarsLiteralData()
	r.FormatVarsJSONData()
	r.FormatVars()
	r.FormatStructs()
	r.FormatStorageTypes()
	r.FormatStorageVars()
	r.FormatStorageStorages()
	r.FormatStorageLinks()
	r.FormatStorageCategories()
	r.FormatStorage()
}

func (r *RenderGo) FormatStructs() {
	formatter := NewFormatterGoStructs(c.identifier)
	formatter.FormatStruct()
	formatter.FormatStructEqual()
	content := formatter.Close()
	filePath := filepath.Join(path.ExportAbsPath(), "structs", "structs.go")
	c.contentMap[filePath] = content
}

func (r *RenderGo) FormatStorageVars() {
	formatter := NewFormatterGoStorageDynamics()
	formatter.FormatPackages()
	formatter.FormatVars()
	formatter.FormatFuncs()
	formatter.FormatLoading()
	content := formatter.Close()
	filePath := filepath.Join(path.ExportAbsPath(), "storage", "dynamics.go")
	c.contentMap[filePath] = content
}

func (r *RenderGo) FormatStorageStorages() {
	formatter := NewFormatterGoStorageStatics(c.collection.PackageNames(), c.collection.Storages())
	formatter.FormatPackages()
	formatter.FormatVars()
	formatter.FormatFuncs()
	formatter.FormatLoading()
	content := formatter.Close()
	filePath := filepath.Join(path.ExportAbsPath(), "storage", "statics.go")
	c.contentMap[filePath] = content
}

func (r *RenderGo) FormatStorageLinks() {
	formatter := NewFormatterGoStorageLinks(c.collection.Links())
	formatter.FormatPackages()
	formatter.FormatVars()
	formatter.FormatFuncs()
	formatter.FormatLoading()
	content := formatter.Close()
	filePath := filepath.Join(path.ExportAbsPath(), "storage", "links.go")
	c.contentMap[filePath] = content
}

func (r *RenderGo) FormatStorageCategories() {
	formatter := NewFormatterGoStorageCategories(c.collection.Categories())
	formatter.FormatPackages()
	formatter.FormatFuncs()
	formatter.FormatLoading()
	content := formatter.Close()
	filePath := filepath.Join(path.ExportAbsPath(), "storage", "categories.go")
	c.contentMap[filePath] = content
}

func (r *RenderGo) FormatStorage() {
	formatter := NewFormatterGoStorage()
	formatter.FormatPackages()
	formatter.FormatVars()
	formatter.FormatFuncs()
	formatter.FormatLoading()
	content := formatter.Close()
	filePath := filepath.Join(path.ExportAbsPath(), "storage", "storage.go")
	c.contentMap[filePath] = content
}

func (r *RenderGo) FormatStorageTypes() {
	formatter := NewFormatterGoStorageTypes(c.identifier)
	formatter.FormatPackages()
	formatter.FormatTypes()
	formatter.FormatVars()
	formatter.FormatFuncs()
	var nodes = []Node{}
	for _, domain := range c.excelMap[FlagBase] {
		for _, excel := range domain[ExcelTypeRegular] {
			for _, node := range excel.Nodes() {
				if c.FilterNodeByDataType(node) {
					nodes = append(nodes, node)
				}
			}
		}
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].String() < nodes[j].String()
	})
	for _, node := range nodes {
		formatter.FormatNode(node)
	}
	content := formatter.Close()
	filePath := filepath.Join(path.ExportAbsPath(), "storage", "types.go")
	c.contentMap[filePath] = content
}

func (r *RenderGo) FormatVarsLiteralData() {
	domains := make([]Domain, 0)
	c.ForeachDomain(func(domain Domain) {
		domains = append(domains, domain)
	})
	results := c.Parallel(ToSlice(domains), func(param any) func() any {
		return func() any {
			domain := param.(Domain)
			formatter := NewFormatterGoLiteralData(r.GetPackageName(domain), c.identifier)
			for _, excel := range domain[ExcelTypeRegular] {
				for _, node := range excel.Nodes() {
					if c.FilterNodeByDataType(node) {
						formatter.FormatNode(node)
					}
				}
			}
			content := formatter.Close()
			if len(content) == 0 {
				return nil
			}
			return []string{r.GetLiteralDataPath(domain), content}
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

func (r *RenderGo) FormatVarsJSONData() {
	domains := make([]Domain, 0)
	c.ForeachDomain(func(domain Domain) {
		domains = append(domains, domain)
	})
	results := c.Parallel(ToSlice(domains), func(param any) func() any {
		return func() any {
			domain := param.(Domain)
			formatter := NewFormatterGoVarsJSONData(r.GetPackageName(domain), c.identifier)
			for _, excel := range domain[ExcelTypeRegular] {
				for _, node := range excel.Nodes() {
					if c.FilterNodeByDataType(node) {
						formatter.FormatNode(node)
					}
				}
			}
			content := formatter.Close()
			if len(content) == 0 {
				return nil
			}
			return []string{r.GetJSONDataPath(domain), content}
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

func (r *RenderGo) FormatVars() {
	packageNames := make([]string, 0)
	for packageName := range c.excelMap {
		packageNames = append(packageNames, packageName)
	}
	results := c.Parallel(ToSlice(packageNames), func(param any) func() any {
		return func() any {
			packageName := param.(string)
			formatter := NewFormatterGoVar(format.ToGoPackageCase(packageName), c.identifier)
			var nodes = []Node{}
			for _, domain := range c.excelMap[packageName] {
				for _, excel := range domain[ExcelTypeRegular] {
					for _, node := range excel.Nodes() {
						if c.FilterNodeByDataType(node) {
							nodes = append(nodes, node)
						}
					}
				}
			}
			sort.Slice(nodes, func(i, j int) bool {
				return nodes[i].String() < nodes[j].String()
			})
			for _, node := range nodes {
				formatter.FormatNode(node)
			}
			content := formatter.Close()
			if len(content) == 0 {
				return nil
			}
			return []string{r.GetVarPath(c, packageName), content}
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

func (r *RenderGo) GetLiteralDataPath(domain Domain) string {
	for _, excels := range domain {
		for _, excel := range excels {
			goPackageName := format.ToGoPackageCase(excel.PackageName())
			fileName := fmt.Sprintf("literal_%s.go", format.ToLower(excel.DomainName()))
			return filepath.Join(path.ExportAbsPath(), "vars", goPackageName, fileName)
		}
	}
	Exit("[Main] Cannot find excel in domain")
	return ""
}

func (r *RenderGo) GetJSONDataPath(domain Domain) string {
	for _, excels := range domain {
		for _, excel := range excels {
			goPackageName := format.ToGoPackageCase(excel.PackageName())
			fileName := fmt.Sprintf("json_%s.go", format.ToLower(excel.DomainName()))
			return filepath.Join(path.ExportAbsPath(), "vars", goPackageName, fileName)
		}
	}
	Exit("[Main] Cannot find excel in domain")
	return ""
}

func (r *RenderGo) GetVarPath(c *Converter, packageName string) string {
	for _, domain := range c.excelMap[packageName] {
		for _, excels := range domain {
			for _, excel := range excels {
				goPackageName := format.ToGoPackageCase(excel.PackageName())
				return filepath.Join(path.ExportAbsPath(), "vars", goPackageName, "vars.go")
			}
		}
	}

	Exit("[Main] Cannot find excel in package")
	return ""
}

func (r *RenderGo) GetPackageName(domain Domain) string {
	for _, excels := range domain {
		for _, excel := range excels {
			return format.ToGoPackageCase(excel.PackageName())
		}
	}
	Exit("[Main] Cannot find excel in domain")
	return ""
}
