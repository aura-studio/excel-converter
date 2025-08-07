package converter

import (
	"fmt"
	"path/filepath"
	"sort"
)

type ConverterGo struct {
	*ConverterBase
	identifier *Identifier
	collection *Collection
}

func NewConverterGo() *ConverterGo {
	c := &ConverterGo{
		ConverterBase: NewConverterBase(ConverterTypeGo, FieldTypeServer),
		identifier:    NewIdentifier(),
		collection:    NewCollection(),
	}
	return c
}

func (c *ConverterGo) Run() {
	c.Load()
	c.Parse()
	c.Export()
}

func (c *ConverterGo) Parse() {
	c.Build()
	c.Identity()
	c.Link()
}

func (c *ConverterGo) Export() {
	c.Format()
	c.Remove()
	c.Write()
}

func (c *ConverterGo) Identity() {
	nodes := []Node{}
	c.ForeachExcel(func(e Excel) {
		if e.Type() == ExcelTypeTemplate {
			for _, node := range e.Nodes() {
				if node.Excel().ForServer() && node.Sheet().ForServer() {
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
				if node.Excel().ForServer() && node.Sheet().ForServer() {
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
				if node.Excel().ForServer() && node.Sheet().ForServer() {
					c.identifier.GenerateStruct(node)
				}
			}
		}
	})
	c.ForeachExcel(func(e Excel) {
		if e.Type() == ExcelTypeTemplate || e.Type() == ExcelTypeRegular {
			for _, node := range e.Nodes() {
				if node.Excel().ForServer() && node.Sheet().ForServer() {
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

func (c *ConverterGo) Link() {
	c.ForeachExcel(func(e Excel) {
		if e.Type() == ExcelTypeRegular {
			for _, node := range e.Nodes() {
				if node.Excel().ForServer() && node.Sheet().ForServer() {
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

func (c *ConverterGo) Format() {
	c.FormatVarsLiteralData()
	c.FormatVarsJSONData()
	c.FormatVars()
	c.FormatStructs()
	c.FormatStorageTypes()
	c.FormatStorageVars()
	c.FormatStorageStorages()
	c.FormatStorageLinks()
	c.FormatStorageCategories()
	c.FormatStorage()
}

func (c *ConverterGo) FormatStructs() {
	formatter := NewFormatterGoStructs(c.identifier)
	formatter.FormatStruct()
	formatter.FormatStructEqual()
	content := formatter.Close()
	filePath := filepath.Join(path.ExportAbsPath(), "structs", "structs.go")
	c.contentMap[filePath] = content
}

func (c *ConverterGo) FormatStorageVars() {
	formatter := NewFormatterGoStorageDynamics()
	formatter.FormatPackages()
	formatter.FormatVars()
	formatter.FormatFuncs()
	formatter.FormatLoading()
	content := formatter.Close()
	filePath := filepath.Join(path.ExportAbsPath(), "storage", "dynamics.go")
	c.contentMap[filePath] = content
}

func (c *ConverterGo) FormatStorageStorages() {
	formatter := NewFormatterGoStorageStatics(c.collection.PackageNames(), c.collection.Storages())
	formatter.FormatPackages()
	formatter.FormatVars()
	formatter.FormatFuncs()
	formatter.FormatLoading()
	content := formatter.Close()
	filePath := filepath.Join(path.ExportAbsPath(), "storage", "statics.go")
	c.contentMap[filePath] = content
}

func (c *ConverterGo) FormatStorageLinks() {
	formatter := NewFormatterGoStorageLinks(c.collection.Links())
	formatter.FormatPackages()
	formatter.FormatVars()
	formatter.FormatFuncs()
	formatter.FormatLoading()
	content := formatter.Close()
	filePath := filepath.Join(path.ExportAbsPath(), "storage", "links.go")
	c.contentMap[filePath] = content
}

func (c *ConverterGo) FormatStorageCategories() {
	formatter := NewFormatterGoStorageCategories(c.collection.Categories())
	formatter.FormatPackages()
	formatter.FormatFuncs()
	formatter.FormatLoading()
	content := formatter.Close()
	filePath := filepath.Join(path.ExportAbsPath(), "storage", "categories.go")
	c.contentMap[filePath] = content
}

func (c *ConverterGo) FormatStorage() {
	formatter := NewFormatterGoStorage()
	formatter.FormatPackages()
	formatter.FormatVars()
	formatter.FormatFuncs()
	formatter.FormatLoading()
	content := formatter.Close()
	filePath := filepath.Join(path.ExportAbsPath(), "storage", "storage.go")
	c.contentMap[filePath] = content
}

func (c *ConverterGo) FormatStorageTypes() {
	formatter := NewFormatterGoStorageTypes(c.identifier)
	formatter.FormatPackages()
	formatter.FormatTypes()
	formatter.FormatVars()
	formatter.FormatFuncs()
	var nodes = []Node{}
	for _, domain := range c.excelMap[FlagBase] {
		for _, excel := range domain[ExcelTypeRegular] {
			for _, node := range excel.Nodes() {
				if node.Excel().ForServer() && node.Sheet().ForServer() {
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

func (c *ConverterGo) FormatVarsLiteralData() {
	domains := make([]Domain, 0)
	c.ForeachDomain(func(domain Domain) {
		domains = append(domains, domain)
	})
	results := c.Parallel(ToSlice(domains), func(param any) func() any {
		return func() any {
			domain := param.(Domain)
			formatter := NewFormatterGoLiteralData(c.GetPackageName(domain), c.identifier)
			for _, excel := range domain[ExcelTypeRegular] {
				for _, node := range excel.Nodes() {
					if node.Excel().ForServer() && node.Sheet().ForServer() {
						formatter.FormatNode(node)
					}
				}
			}
			content := formatter.Close()
			if len(content) == 0 {
				return nil
			}
			return []string{c.GetLiteralDataPath(domain), content}
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

func (c *ConverterGo) FormatVarsJSONData() {
	domains := make([]Domain, 0)
	c.ForeachDomain(func(domain Domain) {
		domains = append(domains, domain)
	})
	results := c.Parallel(ToSlice(domains), func(param any) func() any {
		return func() any {
			domain := param.(Domain)
			formatter := NewFormatterGoVarsJSONData(c.GetPackageName(domain), c.identifier)
			for _, excel := range domain[ExcelTypeRegular] {
				for _, node := range excel.Nodes() {
					if node.Excel().ForServer() && node.Sheet().ForServer() {
						formatter.FormatNode(node)
					}
				}
			}
			content := formatter.Close()
			if len(content) == 0 {
				return nil
			}
			return []string{c.GetJSONDataPath(domain), content}
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

func (c *ConverterGo) FormatVars() {
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
						if node.Excel().ForServer() && node.Sheet().ForServer() {
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
			return []string{c.GetVarPath(packageName), content}
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

func (c *ConverterGo) GetLiteralDataPath(domain Domain) string {
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

func (c *ConverterGo) GetJSONDataPath(domain Domain) string {
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

func (c *ConverterGo) GetVarPath(packageName string) string {
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

func (c *ConverterGo) GetPackageName(domain Domain) string {
	for _, excels := range domain {
		for _, excel := range excels {
			return format.ToGoPackageCase(excel.PackageName())
		}
	}
	Exit("[Main] Cannot find excel in domain")
	return ""
}
