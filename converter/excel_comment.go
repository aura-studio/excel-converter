package converter

type ExcelComment struct {
	*ExcelBase
}

func NewExcelComment(path Path, relPath string, fieldType FieldType) *ExcelComment {
	return &ExcelComment{
		ExcelBase: NewExcelBase(path, relPath, fieldType),
	}
}

func (*ExcelComment) Type() ExcelType {
	return ExcelTypeComment
}
