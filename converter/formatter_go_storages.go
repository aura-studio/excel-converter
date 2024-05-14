package converter

type FormatterGoStorages struct {
	*FormatterBase
	exelMap map[string]map[string]bool
}

func NewFormatterGoStorages() *FormatterGoStorages {
	f := &FormatterGoStorages{
		FormatterBase: NewFormatterBase(),
		exelMap:       make(map[string]map[string]bool),
	}
	f.WriteString(`// <important: auto generate by excel-to-go converter, do not modify>
package storage

`)
	return f
}

func (f *FormatterGoStorages) Close() string {
	return f.String()
}

func (f *FormatterGoStorages) FormatPackage(packageNames []string) {
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
}

func (f *FormatterGoStorages) FormatVars() {
	f.WriteString(`
var Storage = make(map[string]map[string]map[string]any)
var OriginStorage = make(map[string]map[string]map[string]any)
`)
}

func (f *FormatterGoStorages) FormatFuncs() {
	f.WriteString(`
func Set(packageName, excelName, sheetName string, v any) {
	if _, ok := Storage[packageName]; !ok {
		Storage[packageName] = make(map[string]map[string]any)
		OriginStorage[packageName] = make(map[string]map[string]any)
	}
	if _, ok := Storage[packageName][excelName]; !ok {
		Storage[packageName][excelName] = make(map[string]any)
		OriginStorage[packageName][excelName] = make(map[string]any)
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

func Has(packageName, excelName, sheetName string) bool {
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

func Find(packageName, excelName, sheetName string) any {
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

func Link(dstCategory, dstExcelName, dstSheetName,srcCategory, srcExcelName, srcSheetName string) {
	Set(dstCategory, dstExcelName, dstSheetName, Find(srcCategory, srcExcelName, srcSheetName))
}

func Load(dataMap map[string]string, name string, v any) {
	if _, ok := dataMap[name]; !ok {
		return
	}
	if err := json.Unmarshal([]byte(dataMap[name]), v); err != nil {
		panic(err)
	}
}
`)
}

func (f *FormatterGoStorages) FormatStorages(storages []*Storage) {
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

func (f *FormatterGoStorages) FormatLinks(links []*Link) {
	f.WriteString(`
func LinksMapping() {
`)
	for _, link := range links {
		f.WriteString("\tLink(\"")
		f.WriteString(link.DstLinkPath.Category)
		f.WriteString("\", \"")
		f.WriteString(link.DstLinkPath.ExcelName)
		f.WriteString("\", \"")
		f.WriteString(link.DstLinkPath.SheetName)
		f.WriteString("\", \"")
		f.WriteString(link.SrcLinkPath.Category)
		f.WriteString("\", \"")
		f.WriteString(link.SrcLinkPath.ExcelName)
		f.WriteString("\", \"")
		f.WriteString(link.SrcLinkPath.SheetName)
		f.WriteString("\")\n")
	}
	f.WriteString("}\n")
}

func (f *FormatterGoStorages) FormatCategoryLinks() {
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
