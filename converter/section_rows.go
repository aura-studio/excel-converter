package converter

type SectionRows struct {
	*SectionBase
}

func NewSectionRows(sheet Sheet) *SectionRows {
	return &SectionRows{
		SectionBase: NewSectionBase(sheet),
	}
}

// Preprocess 切成横向条形
func (s *SectionRows) Format(int) {
	s.TrimCells()
	s.TrimColumns()
	s.CutRows()
}

func (s *SectionRows) Type() SectionType {
	return SectionTypeRows
}
