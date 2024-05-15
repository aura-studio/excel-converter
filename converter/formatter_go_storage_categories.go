package converter

type FormatterGoStorageCategories struct {
	*FormatterBase
	categories []string
}

func NewFormatterGoStorageCategories(categories []string) *FormatterGoStorageCategories {
	f := &FormatterGoStorageCategories{
		FormatterBase: NewFormatterBase(),
		categories:    categories,
	}
	f.WriteString(`// <important: auto generate by excel-to-go converter, do not modify>
package storage
`)
	return f
}

func (f *FormatterGoStorageCategories) FormatPackages() {
	// Do nothing
}

func (f *FormatterGoStorageCategories) FormatVars() {
	// Do nothing
}

func (f *FormatterGoStorageCategories) FormatFuncs() {
	f.WriteString(`
func LoadCategory(category string) {
	dstCategory := category
	srcCategory := Parent(dstCategory)
	if srcCategory == "" {
		return
	}
	if _, ok := Storage[dstCategory]; !ok {
		Storage[dstCategory] = Storage[srcCategory]
		return
	}
	for excelName, excel := range Storage[srcCategory] {
		if _, ok := Storage[dstCategory][excelName]; !ok {
			Storage[dstCategory][excelName] = excel
		} else {
			for sheetName, sheet := range Storage[srcCategory][excelName] {
				if _, ok := Storage[dstCategory][excelName][sheetName]; !ok {
					Storage[dstCategory][excelName][sheetName] = sheet
				}
			}
		}
	}
}
`)
}

func (f *FormatterGoStorageCategories) FormatLoading() {
	f.WriteString(`
func LoadCategories() {
`)
	for _, category := range f.categories {
		f.WriteString("\tLoadCategory(\"")
		f.WriteString(category)
		f.WriteString("\")\n")
	}
	f.WriteString("}\n")
}

func (f *FormatterGoStorageCategories) Close() string {
	return f.String()
}
