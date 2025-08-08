package converter

type ExcelTemplate struct {
	*ExcelRegular
}

func NewExcelTemplate(path Path, relPath string) *ExcelTemplate {
	return &ExcelTemplate{
		ExcelRegular: NewExcelRegular(path, relPath),
	}
}

func (*ExcelTemplate) Type() ExcelType {
	return ExcelTypeTemplate
}
