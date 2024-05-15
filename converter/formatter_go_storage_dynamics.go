package converter

type FormatterGoStorageDynamics struct {
	*FormatterBase
	packageNames []string
	storages     []*Storage
}

func NewFormatterGoStorageDynamics(packageNames []string, storages []*Storage) *FormatterGoStorageDynamics {
	f := &FormatterGoStorageDynamics{
		FormatterBase: NewFormatterBase(),
		packageNames:  packageNames,
		storages:      storages,
	}
	f.WriteString(`// <important: auto generate by excel-to-go converter, do not modify>
package storage
`)
	return f
}

func (f *FormatterGoStorageDynamics) FormatPackages() {
	f.WriteString(`
import (
	"encoding/json"
	"strings"

	"github.com/mohae/deepcopy"
)
`)
}

func (f *FormatterGoStorageDynamics) FormatVars() {
	// Do nothing
}

func (f *FormatterGoStorageDynamics) FormatFuncs() {
	f.WriteString(`
func LoadDynamic(data map[string]string, packageName, excelName, sheetName string) {
	if data == nil {
		return
	}
	name := strings.Join([]string{packageName, excelName, sheetName}, ".")
	if _, ok := data[name]; !ok {
		return
	}
	excelTypeStorages, ok := TypeStorage[excelName]
	if !ok {
		return
	}
	sheetTypeStorage, ok := excelTypeStorages[sheetName]
	if !ok {
		return
	}
	v := deepcopy.Copy(sheetTypeStorage)
	if err := json.Unmarshal([]byte(data[name]), &v); err != nil {
		panic(err)
	}
	if _, ok := OriginStorage[packageName]; !ok {
		OriginStorage[packageName] = make(map[string]map[string]any)
	}
	if _, ok := OriginStorage[packageName][excelName]; !ok {
		OriginStorage[packageName][excelName] = make(map[string]any)
	}
	OriginStorage[packageName][excelName][sheetName] = v
}

`)
}

func (f *FormatterGoStorageDynamics) FormatLoading() {
	f.WriteString(`
func LoadDynamics(data map[string]string) {
`)
	for _, storage := range f.storages {
		f.WriteString("\tLoadDynamic(data, \"")
		f.WriteString(storage.StoragePath.PackageName)
		f.WriteString("\", \"")
		f.WriteString(storage.StoragePath.ExcelName)
		f.WriteString("\", \"")
		f.WriteString(storage.StoragePath.SheetName)
		f.WriteString("\")\n")
	}
	f.WriteString("}\n")
}

func (f *FormatterGoStorageDynamics) Close() string {
	return f.String()
}
