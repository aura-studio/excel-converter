package converter

type FormatterGoLinks struct {
	*FormatterBase
	links []*Link
}

func NewFormatterGoStorageLinks(links []*Link) *FormatterGoLinks {
	f := &FormatterGoLinks{
		FormatterBase: NewFormatterBase(),
		links:         links,
	}
	f.WriteString(`// <important: auto generate by excel-to-go converter, do not modify>
package storage
`)
	return f
}

func (f *FormatterGoLinks) FormatPackages() {
	// Do nothing
}

func (f *FormatterGoLinks) FormatVars() {
	// Do nothing
}

func (f *FormatterGoLinks) FormatFuncs() {
	f.WriteString(`
func LoadLink(dstPackageName, dstExcelName, dstSheetName, srcPackageName, srcExcelName, srcSheetName string) {
	for {
		if srcPackageName == "" {
			return
		}
		if _, ok := Storage[srcPackageName]; !ok {
			srcPackageName = Parent(srcPackageName)
			continue
		}
		if _, ok := Storage[srcPackageName][srcExcelName]; !ok {
			srcPackageName = Parent(srcPackageName)
			continue
		}
		if _, ok := Storage[srcPackageName][srcExcelName][srcSheetName]; !ok {
			srcPackageName = Parent(srcPackageName)
			continue
		}
		v := Storage[srcPackageName][srcExcelName][srcSheetName]
		if _, ok := Storage[dstPackageName]; !ok {
			Storage[dstPackageName] = make(map[string]map[string]any)
		}
		if _, ok := Storage[dstPackageName][dstExcelName]; !ok {
			Storage[dstPackageName][dstExcelName] = make(map[string]any)
		}
		Storage[dstPackageName][dstExcelName][dstSheetName] = v
	}
}
`)
}

func (f *FormatterGoLinks) FormatLoading() {
	f.WriteString(`
func LoadLinks() {
`)
	for _, link := range f.links {
		f.WriteString("\tLoadLink(\"")
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

func (f *FormatterGoLinks) Close() string {
	return f.String()
}
