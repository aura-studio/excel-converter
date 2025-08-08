package converter

type ExcelComment struct {
	*ExcelBase
}

func NewExcelComment(path Path, relPath string) *ExcelComment {
	return &ExcelComment{
		ExcelBase: NewExcelBase(path, relPath),
	}
}

func (*ExcelComment) Type() ExcelType {
	return ExcelTypeComment
}
