package converter

type FormatterGoStorage struct {
	*FormatterBase
}

func NewFormatterGoStorage() *FormatterGoStorage {
	f := &FormatterGoStorage{
		FormatterBase: NewFormatterBase(),
	}
	f.WriteString(`// <important: auto generate by excel-to-go converter, do not modify>
package storage
`)
	return f
}

func (f *FormatterGoStorage) FormatPackages() {
	f.WriteString(`
import (
	"strings"

	"github.com/mohae/deepcopy"
)
`)
}

func (f *FormatterGoStorage) FormatVars() {
	f.WriteString(`
var Storage = make(map[string]map[string]map[string]any)
var OriginStorage = make(map[string]map[string]map[string]any)
`)
}

func (f *FormatterGoStorage) FormatFuncs() {
	f.WriteString(`
func Parent(packageName string) string {
	if packageName == "Base" {
		return ""
	}
	strs := strings.Split(packageName, "_")
	if len(strs) == 1 {
		return "Base"
	} else {
		return strings.Join(strs[:len(strs)-1], "_")
	}
}

func ResetOriginStorage() {
	storage := make(map[string]map[string]map[string]any)
	for packageName, subStorage := range OriginStorage {
		for excelName, excel := range subStorage {
			for sheetName, sheet := range excel {
				if _, ok := storage[packageName]; !ok {
					storage[packageName] = make(map[string]map[string]any)
				}
				if _, ok := storage[packageName][excelName]; !ok {
					storage[packageName][excelName] = make(map[string]any)
				}
				storage[packageName][excelName][sheetName] = deepcopy.Copy(sheet)
			}
		}
	}
	OriginStorage = storage
}

func ResetStorage() {
	storage := make(map[string]map[string]map[string]any)
	for packageName, subStorage := range OriginStorage {
		for excelName, excel := range subStorage {
			for sheetName, sheet := range excel {
				if _, ok := storage[packageName]; !ok {
					storage[packageName] = make(map[string]map[string]any)
				}
				if _, ok := storage[packageName][excelName]; !ok {
					storage[packageName][excelName] = make(map[string]any)
				}
				storage[packageName][excelName][sheetName] = deepcopy.Copy(sheet)
			}
		}
	}
	Storage = storage
}
`)
}

func (f *FormatterGoStorage) FormatLoading() {
	f.WriteString(`
func init() {
	ResetOriginStorage()
	LoadStatics()
	ResetStorage()
	LoadLinks()
	LoadCategories()

	// For dynamic
	LoadTypes()
}

func Load(data map[string]string) {
	ResetOriginStorage()
	LoadDynamics(data)
	ResetStorage()
	LoadLinks()
	LoadCategories()
}
`)
}

func (f *FormatterGoStorage) Close() string {
	return f.String()
}
