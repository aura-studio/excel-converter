package converter

type SectionColumns struct {
	*SectionBase
}

func NewSectionColumns(sheet Sheet) *SectionColumns {
	return &SectionColumns{
		SectionBase: NewSectionBase(sheet),
	}
}

// Preprocess 旋转后，切成横向条形
func (s *SectionColumns) Format(int) {
	s.TrimCells()
	s.Rotate()
	s.TrimColumns()
	s.CutRows()
}

func (s *SectionColumns) Type() SectionType {
	return SectionTypeColumns
}
