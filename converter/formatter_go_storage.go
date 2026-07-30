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
)
`)
}

func (f *FormatterGoStorage) FormatVars() {
	f.WriteString(`
// Storage 是最终生效的配置：已应用 link 与 category 解析。
// 它在 init 阶段一次性构建完成，之后只读。
var Storage = make(map[string]map[string]map[string]any)

// originStorage 是 LoadStatics 直接写入的原始配置，未经 link / category 解析，
// 仅作为构建 Storage 的中间产物，不对外暴露。
var originStorage = make(map[string]map[string]map[string]any)
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

// buildStorage 由 originStorage 重建三层 map 骨架，叶子（sheet）数据按引用共享。
//
// 配置表在运行期是只读的：需要修改表数据的调用方会自行深拷贝再改。
// 因此这里不做 deepcopy，否则 Storage 与 originStorage 会各自持有
// 一整套表数据，启动内存成倍放大。
func buildStorage() map[string]map[string]map[string]any {
	storage := make(map[string]map[string]map[string]any, len(originStorage))
	for packageName, subStorage := range originStorage {
		for excelName, excel := range subStorage {
			for sheetName, sheet := range excel {
				if _, ok := storage[packageName]; !ok {
					storage[packageName] = make(map[string]map[string]any, len(subStorage))
				}
				if _, ok := storage[packageName][excelName]; !ok {
					storage[packageName][excelName] = make(map[string]any, len(excel))
				}
				storage[packageName][excelName][sheetName] = sheet
			}
		}
	}
	return storage
}
`)
}

func (f *FormatterGoStorage) FormatLoading() {
	f.WriteString(`
func init() {
	LoadStatics()
	Storage = buildStorage()
	LoadLinks()
	LoadCategories()
}
`)
}

func (f *FormatterGoStorage) Close() string {
	return f.String()
}
