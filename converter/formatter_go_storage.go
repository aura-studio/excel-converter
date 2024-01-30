package converter

type FormatterGoStorage struct {
	*FormatterBase
	exelMap map[string]map[string]bool
}

func NewFormatterGoStorage() *FormatterGoStorage {
	f := &FormatterGoStorage{
		FormatterBase: NewFormatterBase(),
		exelMap:       make(map[string]map[string]bool),
	}
	f.WriteString(`// <important: auto generate by excel-to-go converter, do not modify>
package storage

`)
	return f
}

func (f *FormatterGoStorage) Close() string {
	f.WriteString(`
func Set(packageName string, excelName string, sheetName string, v interface{}) {
	if _, ok := Storage[packageName]; !ok {
		Storage[packageName] = make(map[string]map[string]interface{})
		OriginStorage[packageName] = make(map[string]map[string]interface{})
	}
	if _, ok := Storage[packageName][excelName]; !ok {
		Storage[packageName][excelName] = make(map[string]interface{})
		OriginStorage[packageName][excelName] = make(map[string]interface{})
	}
	Storage[packageName][excelName][sheetName] = v
	OriginStorage[packageName][excelName][sheetName] = v
}

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

func Has(packageName string, excelName string, sheetName string) bool {
	if _, ok := Storage[packageName]; !ok {
		return false
	}
	if _, ok := Storage[packageName][excelName]; !ok {
		return false
	}
	if _, ok := Storage[packageName][excelName][sheetName]; !ok {
		return false
	}
	return true
}

func Find(packageName string, excelName string, sheetName string) interface{} {
	for {
		if packageName == "" {
			return nil
		}
		if !Has(packageName, excelName, sheetName) {
			packageName = Parent(packageName)
			continue
		}
		return Storage[packageName][excelName][sheetName]
	}
}

func Link(dstExcelName string, dstSheetName string, srcExcelName string, srcSheetName string) {
	for _, category := range Categories {
		Set(category, dstExcelName, dstSheetName, Find(category, srcExcelName, srcSheetName))
	}
}

func Load(dataMap map[string]string, name string, v interface{}) {
	if _, ok := dataMap[name]; !ok {
		return
	}
	if err := json.Unmarshal([]byte(dataMap[name]), v); err != nil {
		panic(err)
	}
}

`)
	return f.String()
}

func (f *FormatterGoStorage) FormatPackage(packageNames []string) {
	f.WriteString("import (\n")
	f.WriteString("\t\"encoding/json\"\n")
	for _, packageName := range packageNames {
		f.WriteString("\t\"")
		f.WriteString(path.ImportPath() + "/")
		goPackageName := format.ToGoPackageCase(packageName)
		f.WriteString(goPackageName)
		f.WriteString("\"\n")
	}
	f.WriteString("\t\"strings\"\n")
	f.WriteString(")\n")
	f.WriteString(`
var Storage = make(map[string]map[string]map[string]interface{})
var OriginStorage = make(map[string]map[string]map[string]interface{})

`)
}

func (f *FormatterGoStorage) FormatCategories(categories []string) {
	f.WriteString("var Categories = []string{\n")
	for _, category := range categories {
		f.WriteString("\t\"")
		f.WriteString(category)
		f.WriteString("\",\n")
	}
	f.WriteString("}\n")
}

func (f *FormatterGoStorage) FormatStorages(storages []*Storage) {
	f.WriteString(`
func StoragesLoading(data map[string]string) {
`)
	for _, storage := range storages {
		f.WriteString("\tLoad(data, \"")
		f.WriteString(storage.StoragePath.PackageName)
		f.WriteString(".")
		f.WriteString(storage.StoragePath.ExcelName)
		f.WriteString(".")
		f.WriteString(storage.StoragePath.SheetName)
		f.WriteString("\", &")
		f.WriteString(format.ToGoPackageCase(storage.StorageVar.PackageName))
		f.WriteString(".")
		f.WriteString(storage.StorageVar.VarName)
		f.WriteString(")\n")
	}
	f.WriteString("}\n")

	f.WriteString(`
func StoragesMapping() {
`)
	for _, storage := range storages {
		f.WriteString("\tSet(\"")
		f.WriteString(storage.StoragePath.PackageName)
		f.WriteString("\", \"")
		f.WriteString(storage.StoragePath.ExcelName)
		f.WriteString("\", \"")
		f.WriteString(storage.StoragePath.SheetName)
		f.WriteString("\", ")
		f.WriteString(format.ToGoPackageCase(storage.StorageVar.PackageName))
		f.WriteString(".")
		f.WriteString(storage.StorageVar.VarName)
		f.WriteString(")\n")
	}
	f.WriteString("}\n")
}

func (f *FormatterGoStorage) StoragesLoading(storages []*Storage) {

}

func (f *FormatterGoStorage) FormatLinks(links []*Link) {
	f.WriteString(`
func LinksMapping() {
`)
	for _, link := range links {
		f.WriteString("\tLink(\"")
		f.WriteString(link.DstLinkPath.ExcelName)
		f.WriteString("\", \"")
		f.WriteString(link.DstLinkPath.SheetName)
		f.WriteString("\", \"")
		f.WriteString(link.SrcLinkPath.ExcelName)
		f.WriteString("\", \"")
		f.WriteString(link.SrcLinkPath.SheetName)
		f.WriteString("\")\n")
	}
	f.WriteString("}\n")
}

func (f *FormatterGoStorage) FormatCategoryLinks() {
	f.WriteString(`
func CategoriesMapping() {
	for _, dstCategory := range Categories {
		srcCategory := Parent(dstCategory)
		if srcCategory == "" {
			continue
		}
		if _, ok := Storage[dstCategory]; !ok {
			Storage[dstCategory] = Storage[srcCategory]
		} else {
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
	}
}
`)
}
