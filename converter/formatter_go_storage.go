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

func CopyStorage() {
	Storage = deepcopy.Copy(OriginStorage).(map[string]map[string]map[string]any)
}
`)
}

func (f *FormatterGoStorage) FormatLoading() {
	f.WriteString(`
func init() {
	LoadStatics()
	CopyStorage()
	LoadLinks()
	LoadCategories()

	// For dynamic
	LoadTypes()
}

func Load(data map[string]string) {
	LoadDynamics(data)
	CopyStorage()
	LoadLinks()
	LoadCategories()
}
`)
}

func (f *FormatterGoStorage) Close() string {
	return f.String()
}
