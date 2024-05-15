package converter

type FormatterGoStorageStatics struct {
	*FormatterBase
	packageNames []string
	storages     []*Storage
}

func NewFormatterGoStorageStatics(packageNames []string, storages []*Storage) *FormatterGoStorageStatics {
	f := &FormatterGoStorageStatics{
		FormatterBase: NewFormatterBase(),
		packageNames:  packageNames,
		storages:      storages,
	}
	f.WriteString(`// <important: auto generate by excel-to-go converter, do not modify>
package storage
`)
	return f
}

func (f *FormatterGoStorageStatics) FormatPackages() {
	f.WriteString("\nimport (\n")
	for _, packageName := range f.packageNames {
		f.WriteString("\t\"")
		f.WriteString(path.ImportPath())
		f.WriteString("/vars/")
		goPackageName := format.ToGoPackageCase(packageName)
		f.WriteString(goPackageName)
		f.WriteString("\"\n")
	}
	f.WriteString(")\n")
}

func (f *FormatterGoStorageStatics) FormatVars() {
	// Do nothing
}

func (f *FormatterGoStorageStatics) FormatFuncs() {
	f.WriteString(`
func LoadStatic(packageName, excelName, sheetName string, v any) {
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

func (f *FormatterGoStorageStatics) FormatLoading() {
	f.WriteString(`
func LoadStatics() {
`)
	for _, storage := range f.storages {
		f.WriteString("\tLoadStatic(\"")
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

func (f *FormatterGoStorageStatics) Close() string {
	return f.String()
}
