package converter

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
)

type Domain map[ExcelType][]Excel

type Task func() any

const (
	maxExcelReadWorkers  = 1024
	maxFileWriteWorkers  = 256
	fileWriteWorkerScale = 16
)

type Converter struct {
	excelMap      map[string]map[string]Domain
	contentMap    map[string]string
	identifier    *Identifier
	collection    *Collection
	writeDirCache sync.Map
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
	return c.ParallelLimit(params, 0, generator)
}

func (c *Converter) ParallelLimit(
	params []any,
	limit int,
	generator func(any) func() any,
) (results []any) {
	tasks := make([]Task, 0, len(params))
	for _, param := range params {
		tasks = append(tasks, generator(param))
	}

	if limit <= 0 || limit > len(tasks) {
		limit = len(tasks)
	}
	if limit == 0 {
		return nil
	}
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	taskCh := make(chan Task)
	for i := 0; i < limit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {
				result := task()
				mu.Lock()
				results = append(results, result)
				mu.Unlock()
			}
		}()
	}
	for _, task := range tasks {
		taskCh <- task
	}
	close(taskCh)
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
	// ValidateBaseExcelCoverage removed: Base is no longer required to be a superset of all category tables
	c.Read()
	c.Preprocess()
}

func (c *Converter) ValidateBaseExcelCoverage() {
	basePackage, ok := c.excelMap[FlagBase]
	if !ok {
		Exit("[ValidateBaseExcelCoverage] package %s not found", FlagBase)
	}

	baseTables := make(map[string]struct{})
	for _, domain := range basePackage {
		for typ, excels := range domain {
			if typ != ExcelTypeRegular && typ != ExcelTypeTemplate {
				continue
			}
			for _, excel := range excels {
				baseTables[excel.FixedName()] = struct{}{}
			}
		}
	}

	missingTables := make(map[string]map[string]struct{})
	for packageName, pkg := range c.excelMap {
		if packageName == FlagBase || packageName == FlagDefault {
			continue
		}
		for _, domain := range pkg {
			for typ, excels := range domain {
				if typ != ExcelTypeRegular && typ != ExcelTypeTemplate {
					continue
				}
				for _, excel := range excels {
					tableName := excel.FixedName()
					if _, exists := baseTables[tableName]; exists {
						continue
					}
					if _, exists := missingTables[tableName]; !exists {
						missingTables[tableName] = make(map[string]struct{})
					}
					missingTables[tableName][packageName] = struct{}{}
				}
			}
		}
	}

	if len(missingTables) == 0 {
		return
	}

	tableNames := make([]string, 0, len(missingTables))
	for tableName := range missingTables {
		tableNames = append(tableNames, tableName)
	}
	sort.Strings(tableNames)

	var builder strings.Builder
	builder.WriteString("\n")
	builder.WriteString(strings.Repeat("=", 72))
	builder.WriteString("\n")
	builder.WriteString("VALIDATION FAILED: BASE TABLE COVERAGE\n")
	builder.WriteString(strings.Repeat("=", 72))
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("Missing table count: %d\n\n", len(tableNames)))
	for idx, tableName := range tableNames {
		categories := make([]string, 0, len(missingTables[tableName]))
		for category := range missingTables[tableName] {
			categories = append(categories, category)
		}
		sort.Strings(categories)
		builder.WriteString(fmt.Sprintf("%d) TABLE: %s\n", idx+1, tableName))
		builder.WriteString(fmt.Sprintf("   Categories: %s\n\n", strings.Join(categories, ", ")))
	}
	builder.WriteString("Action required: add the tables above to Base, then rerun converter.\n")
	builder.WriteString(strings.Repeat("=", 72))

	Exit("[ValidateBaseExcelCoverage] %s", builder.String())
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
				buf := new(bytes.Buffer)
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
	c.ParallelLimit(ToSlice(excels), excelReadWorkers(), func(param any) func() any {
		return func() any {
			excel := param.(Excel)
			excel.Read()
			return nil
		}
	})
}

func excelReadWorkers() int {
	workers := runtime.NumCPU() * 128
	if workers > maxExcelReadWorkers {
		workers = maxExcelReadWorkers
	}
	if workers < 1 {
		workers = 1
	}
	return workers
}

func (c *Converter) Write() {
	absPaths := make([]any, 0, len(c.contentMap))
	for absPath := range c.contentMap {
		absPaths = append(absPaths, absPath)
	}
	c.ParallelLimit(absPaths, fileWriteWorkers(), func(param any) func() any {
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

func fileWriteWorkers() int {
	workers := runtime.NumCPU() * fileWriteWorkerScale
	if workers > maxFileWriteWorkers {
		workers = maxFileWriteWorkers
	}
	if workers < 1 {
		workers = 1
	}
	return workers
}

func (c *Converter) WriteFile(absPath string, s string) error {
	fileDir := filepath.Dir(absPath)
	if _, ok := c.writeDirCache.Load(fileDir); !ok {
		if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
			Exit(fmt.Errorf("[%s], %v", absPath, err))
		}
		c.writeDirCache.Store(fileDir, struct{}{})
	}
	data := []byte(s)
	if oldData, err := os.ReadFile(absPath); err == nil && bytes.Equal(oldData, data) {
		return nil
	}
	return os.WriteFile(absPath, data, 0666)
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
	exportPath := path.ExportAbsPath()
	if err := os.MkdirAll(exportPath, os.ModePerm); err != nil {
		Exit("[Main] Mkdir error, %v", err)
	}

	keep := make(map[string]struct{}, len(c.contentMap))
	for absPath := range c.contentMap {
		keep[filepath.Clean(absPath)] = struct{}{}
	}

	if err := filepath.WalkDir(exportPath, func(absPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		absPath = filepath.Clean(absPath)
		if _, ok := keep[absPath]; ok {
			return nil
		}
		if isGeneratedFile(absPath) {
			return os.Remove(absPath)
		}
		return nil
	}); err != nil {
		Exit("[Main] Remove stale file error, %v", err)
	}
}

func isGeneratedFile(absPath string) bool {
	switch filepath.Ext(absPath) {
	case ".go", ".json", ".lua", ".cs":
	default:
		return false
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		return false
	}
	return bytes.Contains(data, []byte("auto generate by excel-to-"))
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
				if c.FilterNodeByDataType(node) && !skipVectorDataNode(node) {
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
