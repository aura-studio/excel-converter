package converter

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type Domain map[ExcelType][]Excel

type Task func() any

type Converter struct {
	excelMap   map[string]map[string]Domain
	contentMap map[string]string
	identifier *Identifier
	collection *Collection
}

var c = NewConverter()

func NewConverter() *Converter {
	return &Converter{
		excelMap:   make(map[string]map[string]Domain),
		contentMap: make(map[string]string),
		identifier: NewIdentifier(),
		collection: NewCollection(),
	}
}

func (c *Converter) Run() {
	c.Load()
	c.Build()
	c.Identity()
	c.Link()
	c.Render()
	c.Remove()
	c.Write()
}

func (c *Converter) Render() {
	if render, ok := renderMap[env.RenderType]; ok {
		render.Render()
	} else {
		Exit(fmt.Errorf("render type %s not found", env.RenderType))
	}
}

func (c *Converter) Parallel(
	params []any,
	generator func(any) func() any,
) (results []any) {
	var tasks = make([]Task, 0, len(params))
	for _, param := range params {
		tasks = append(tasks, generator(param))
	}

	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	for _, task := range tasks {
		wg.Add(1)
		go func(task Task) {
			defer wg.Done()
			result := task()
			mu.Lock()
			defer mu.Unlock()
			results = append(results, result)
		}(task)
	}
	wg.Wait()

	return
}

func (c *Converter) ForeachDomain(f func(Domain)) {
	for _, pkg := range c.excelMap {
		for _, domain := range pkg {
			f(domain)
		}
	}
}
func (c *Converter) ForeachExcel(f func(Excel)) {
	for _, packageName := range c.excelMap {
		for _, domain := range packageName {
			for _, excels := range domain {
				for _, excel := range excels {
					f(excel)
				}
			}
		}
	}
}

func (c *Converter) Load() {
	c.Scan()
	c.Read()
	c.Preprocess()
}

func (c *Converter) Scan() {
	if err := filepath.Walk(path.ImportAbsPath(), func(absPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		fileName := filepath.Base(absPath)
		if fileName[0] == '~' {
			return nil
		}
		if filepath.Ext(fileName) != FlagExt {
			return nil
		}
		relPath, err := filepath.Rel(path.ImportAbsPath(), absPath)
		if err != nil {
			return err
		}
		excelType := c.ExcelType(relPath)
		excel := excelCreators[excelType](path, relPath)
		packageName := excel.PackageName()
		if _, ok := c.excelMap[packageName]; !ok {
			c.excelMap[packageName] = make(map[string]Domain)
		}
		domain := excel.DomainName()
		if _, ok := c.excelMap[packageName][domain]; !ok {
			c.excelMap[packageName][domain] = make(map[ExcelType][]Excel)
		}
		Debug("Excel %v/%v/%v/%v", packageName, domain, excelType, excel.FixedName())
		c.excelMap[packageName][domain][excelType] = append(c.excelMap[packageName][domain][excelType], excel)
		return nil
	}); err != nil {
		Exit("[Main] Scan file error, %v", err)
	}
	c.ForeachDomain(func(domain Domain) {
		for typ, excels := range domain {
			sort.Slice(excels, func(i, j int) bool {
				return excels[i].FixedName() < excels[j].FixedName()
			})
			domain[typ] = excels
		}
	})

	for packageName, pkgExcelMap := range c.excelMap {
		for domain, domainExcelMap := range pkgExcelMap {
			for typ, typeExcels := range domainExcelMap {
				var buf = new(bytes.Buffer)
				buf.WriteString(`[`)
				for index, excel := range typeExcels {
					buf.WriteString(excel.IndirectName())
					if index != len(typeExcels)-1 {
						buf.WriteString(`, `)
					}
				}
				buf.WriteString(`]`)
				Debug("[%v/%v/%v/...] scanned %v", packageName, domain, typ, buf.String())
			}
		}
	}
}

func (c *Converter) Read() {
	excels := make([]Excel, 0)
	c.ForeachExcel(func(excel Excel) {
		excels = append(excels, excel)
	})
	c.Parallel(ToSlice(excels), func(param any) func() any {
		return func() any {
			excel := param.(Excel)
			excel.Read()
			return nil
		}
	})
}

func (c *Converter) Write() {
	var absPaths = make([]any, 0, len(c.contentMap))
	for absPath := range c.contentMap {
		absPaths = append(absPaths, absPath)
	}
	c.Parallel(absPaths, func(param any) func() any {
		return func() any {
			absPath := param.(string)
			content := c.contentMap[absPath]
			Debug("[%v] write %d bytes", absPath, len(content))
			if err := c.WriteFile(absPath, content); err != nil {
				Exit("[%v] Write file error: %v", absPath, err)
			}
			return nil
		}
	})
}

func (c *Converter) WriteFile(absPath string, s string) error {
	// if dir not exists, then create it
	fileDir := filepath.Dir(absPath)
	if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
		Exit(fmt.Errorf("[%s], %v", absPath, err))
	}
	// if already exists then remove it
	if _, err := os.Stat(absPath); err == nil {
		os.Remove(absPath)
	}
	file, err := os.Create(absPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, s)
	if err != nil {
		return err
	}
	return file.Sync()
}

func (c *Converter) ExcelType(path string) ExcelType {
	switch {
	case strings.Contains(path, FlagComment):
		return ExcelTypeComment
	case strings.Contains(path, FlagSettings):
		return ExcelTypeSettings
	case strings.Contains(path, FlagTemplate):
		return ExcelTypeTemplate
	default:
		return ExcelTypeRegular
	}
}

func (c *Converter) Preprocess() {
	excels := make([]Excel, 0)
	c.ForeachExcel(func(excel Excel) {
		excels = append(excels, excel)
	})
	c.Parallel(ToSlice(excels), func(param any) func() any {
		return func() any {
			excel := param.(Excel)
			excel.Preprocess()
			return nil
		}
	})
}

func (c *Converter) Parse() {

}

func (c *Converter) Build() {
	excels := make([]Excel, 0)
	c.ForeachExcel(func(excel Excel) {
		if excel.Type() == ExcelTypeTemplate || excel.Type() == ExcelTypeRegular {
			excels = append(excels, excel)
		}
	})
	c.Parallel(ToSlice(excels), func(param any) func() any {
		return func() any {
			excel := param.(Excel)
			excel.Build()
			return nil
		}
	})
}

func (c *Converter) Remove() {
	err := os.RemoveAll(path.ExportAbsPath())
	if err != nil {
		Exit("[Main] Remove error, %v", err)
	}
	err = os.Mkdir(path.ExportAbsPath(), os.ModePerm)
	if err != nil {
		Exit("[Main] Mkdir error, %v", err)
	}
}

func (c *Converter) Identity() {
	nodes := []Node{}
	c.ForeachExcel(func(e Excel) {
		if e.Type() == ExcelTypeTemplate {
			for _, node := range e.Nodes() {
				if c.FilterNodeByDataType(node) {
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
				if c.FilterNodeByDataType(node) {
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
				if c.FilterNodeByDataType(node) {
					c.identifier.GenerateStruct(node)
				}
			}
		}
	})
	c.ForeachExcel(func(e Excel) {
		if e.Type() == ExcelTypeTemplate || e.Type() == ExcelTypeRegular {
			for _, node := range e.Nodes() {
				if c.FilterNodeByDataType(node) {
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

func (c *Converter) FilterNodeByDataType(node Node) bool {
	switch env.DataType {
	case DataTypeServer:
		return node.Excel().ForServer() && node.Sheet().ForServer()
	case DataTypeClient:
		return node.Excel().ForClient() && node.Sheet().ForClient()
	default:
		Exit(fmt.Errorf("filter node by data type get invalid data type %s", env.DataType))
	}

	return false
}

func (c *Converter) Link() {
	c.ForeachExcel(func(e Excel) {
		if e.Type() == ExcelTypeRegular {
			for _, node := range e.Nodes() {
				if c.FilterNodeByDataType(node) {
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
