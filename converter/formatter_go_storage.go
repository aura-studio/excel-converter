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

func LoadStorage() {
	Storage = make(map[string]map[string]map[string]any)
	for packageName, subStorage := range OriginStorage {
		for excelName, excel := range subStorage {
			for sheetName, sheet := range excel {
				if _, ok := Storage[packageName]; !ok {
					Storage[packageName] = make(map[string]map[string]any)
				}
				if _, ok := Storage[packageName][excelName]; !ok {
					Storage[packageName][excelName] = make(map[string]any)
				}
				Storage[packageName][excelName][sheetName] = deepcopy.Copy(sheet)
			}
		}
	}
}
`)
}

func (f *FormatterGoStorage) FormatLoading() {
	f.WriteString(`
func init() {
	LoadStatics()
	LoadStorage()
	LoadLinks()
	LoadCategories()

	// For dynamic
	LoadTypes()
}

func Load(data map[string]string) {
	LoadDynamics(data)
	LoadStorage()
	LoadLinks()
	LoadCategories()
}
`)
}

func (f *FormatterGoStorage) Close() string {
	return f.String()
}
