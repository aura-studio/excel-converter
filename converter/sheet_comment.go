package converter

type SheetComment struct {
	*SheetBase
}

func NewSheetComment(excel Excel, name string, data [][]string) *SheetComment {
	return &SheetComment{
		SheetBase: NewSheetBase(excel, name, data),
	}
}

func (*SheetComment) Type() SheetType {
	return SheetTypeComment
}
