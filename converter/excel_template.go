package converter

type ExcelTemplate struct {
	*ExcelRegular
}

func NewExcelTemplate(path Path, relPath string, fieldType FieldType) *ExcelTemplate {
	return &ExcelTemplate{
		ExcelRegular: NewExcelRegular(path, relPath, fieldType),
	}
}

func (*ExcelTemplate) Type() ExcelType {
	return ExcelTypeTemplate
}
