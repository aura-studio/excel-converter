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

type Task func() interface{}

type ConverterBase struct {
	typ        ConverterType
	excelMap   map[string]map[string]Domain
	contentMap map[string]string
}

func NewConverterBase(typ ConverterType) *ConverterBase {
	return &ConverterBase{
		typ:        typ,
		excelMap:   make(map[string]map[string]Domain),
		contentMap: make(map[string]string),
	}
}

func (c *ConverterBase) String() string {
	return string(c.typ)
}

func (c *ConverterBase) Run() {
	Exit("[Main] Invalid call: ConverterBase.Run")
}

func (c *ConverterBase) RelPath() string {
	Exit("[Main] Invalid call: ConverterBase.RelPath")
	return ""
}

func (c *ConverterBase) Parallel(
	params []interface{},
	generator func(interface{}) func() interface{},
) (results []interface{}) {
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

func (c *ConverterBase) ForeachDomain(f func(Domain)) {
	for _, pkg := range c.excelMap {
		for _, domain := range pkg {
			f(domain)
		}
	}
}
func (c *ConverterBase) ForeachExcel(f func(Excel)) {
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

func (c *ConverterBase) Load() {
	c.Scan()
	c.Read()
	c.Preprocess()
}

func (c *ConverterBase) Scan() {
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
		var fieldType FieldType
		switch c.typ {
		case ConverterTypeGo:
			fieldType = FieldTypeServer
		case ConverterTypeLua:
			fieldType = FieldTypeClient
		default:
			Exit("[Main] Unsupported converter type %s", c.typ)
		}
		excel := excelCreators[excelType](path, relPath, fieldType)
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

func (c *ConverterBase) Read() {
	excels := make([]Excel, 0)
	c.ForeachExcel(func(excel Excel) {
		excels = append(excels, excel)
	})
	c.Parallel(ToSlice(excels), func(param interface{}) func() interface{} {
		return func() interface{} {
			excel := param.(Excel)
			excel.Read()
			return nil
		}
	})
}

func (c *ConverterBase) Write() {
	var absPaths = make([]interface{}, 0, len(c.contentMap))
	for absPath := range c.contentMap {
		absPaths = append(absPaths, absPath)
	}
	c.Parallel(absPaths, func(param interface{}) func() interface{} {
		return func() interface{} {
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

func (c *ConverterBase) WriteFile(absPath string, s string) error {
	// if dir not exists, then create it
	fileDir := filepath.Dir(absPath)
	if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
		panic(fmt.Errorf("[%s], %v", absPath, err))
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

func (c *ConverterBase) ExcelType(path string) ExcelType {
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

func (c *ConverterBase) Preprocess() {
	excels := make([]Excel, 0)
	c.ForeachExcel(func(excel Excel) {
		excels = append(excels, excel)
	})
	c.Parallel(ToSlice(excels), func(param interface{}) func() interface{} {
		return func() interface{} {
			excel := param.(Excel)
			excel.Preprocess()
			return nil
		}
	})
}

func (c *ConverterBase) Parse() {

}

func (c *ConverterBase) Build() {
	excels := make([]Excel, 0)
	c.ForeachExcel(func(excel Excel) {
		if excel.Type() == ExcelTypeTemplate || excel.Type() == ExcelTypeRegular {
			excels = append(excels, excel)
		}
	})
	c.Parallel(ToSlice(excels), func(param interface{}) func() interface{} {
		return func() interface{} {
			excel := param.(Excel)
			excel.Build()
			return nil
		}
	})
}

func (c *ConverterBase) Remove() {
	err := os.RemoveAll(path.ExportAbsPath())
	if err != nil {
		Exit("[Main] Remove error, %v", err)
	}
	err = os.Mkdir(path.ExportAbsPath(), os.ModePerm)
	if err != nil {
		Exit("[Main] Mkdir error, %v", err)
	}
}
