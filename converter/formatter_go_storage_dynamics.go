package converter

type FormatterGoStorageDynamics struct {
	*FormatterBase
}

func NewFormatterGoStorageDynamics() *FormatterGoStorageDynamics {
	f := &FormatterGoStorageDynamics{
		FormatterBase: NewFormatterBase(),
	}
	f.WriteString(`// <important: auto generate by excel-to-go converter, do not modify>
package storage
`)
	return f
}

func (f *FormatterGoStorageDynamics) FormatPackages() {
	f.WriteString(`
import (
	"strings"
)
`)
}

func (f *FormatterGoStorageDynamics) FormatVars() {
	// Do nothing
}

func (f *FormatterGoStorageDynamics) FormatFuncs() {
	f.WriteString(`
func LoadDynamic(packageName, excelName, sheetName string, jsonStr string) {
	excelTypeStorages, ok := TypeStorage[excelName]
	if !ok {
		return
	}
	f, ok := excelTypeStorages[sheetName]
	if !ok {
		return
	}
	if _, ok := OriginStorage[packageName]; !ok {
		OriginStorage[packageName] = make(map[string]map[string]any)
	}
	if _, ok := OriginStorage[packageName][excelName]; !ok {
		OriginStorage[packageName][excelName] = make(map[string]any)
	}
	OriginStorage[packageName][excelName][sheetName] = f([]byte(jsonStr))
}
`)
}

func (f *FormatterGoStorageDynamics) FormatLoading() {
	f.WriteString(`
func LoadDynamics(data map[string]string) {
	for name, jsonStr := range data {
		if strings.Count(name, ".") != 2 {
			continue
		}
		names := strings.Split(name, ".")
		LoadDynamic(names[0], names[1], names[2], jsonStr)
	}
}
`)
}

func (f *FormatterGoStorageDynamics) Close() string {
	return f.String()
}
