package converter

type ExcelTemplate struct {
	*ExcelRegular
}

func NewExcelTemplate(path Path, relPath string) *ExcelTemplate {
	return &ExcelTemplate{
		ExcelRegular: NewExcelRegular(path, relPath),
	}
}

func (e *ExcelTemplate) Read() {
	e.read(0)
}

func (*ExcelTemplate) Type() ExcelType {
	return ExcelTypeTemplate
}
