package converter

type FormatterGoCategories struct {
	*FormatterBase
	exelMap map[string]map[string]bool
}

func NewFormatterGoLoading() *FormatterGoCategories {
	f := &FormatterGoCategories{
		FormatterBase: NewFormatterBase(),
		exelMap:       make(map[string]map[string]bool),
	}
	f.WriteString(`// <important: auto generate by excel-to-go converter, do not modify>
package storage

`)
	return f
}

func (f *FormatterGoCategories) FormatVars() {
	f.WriteString(`
var Categories = make([]string, 0)
`)
}

func (f *FormatterGoCategories) FormatFuncs() {
	f.WriteString(`
func Category(category string) {
	Categories = append(Categories, category)
}
`)
}

func (f *FormatterGoStorages) FormatCategories(categories []string) {
	f.WriteString(`
func CategoriesLoading() {
`)
	for _, category := range categories {
		f.WriteString("\tCategory(\"")
		f.WriteString(category)
		f.WriteString("\")\n")
	}
	f.WriteString("}\n")
}

func (f *FormatterGoCategories) Close() string {
	return f.String()
}
