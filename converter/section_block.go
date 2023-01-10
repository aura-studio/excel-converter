package converter

type SectionBlock struct {
	*SectionBase
}

func NewSectionBlock(sheet Sheet) *SectionBlock {
	return &SectionBlock{
		SectionBase: NewSectionBase(sheet),
	}
}

// Preprocess 切成方块形
func (s *SectionBlock) Format(originHeaderSize int) {
	s.TrimCells()
	s.TrimColumns()
	s.CutRows()
	s.FillColumns(originHeaderSize)
}

func (s *SectionBlock) Type() SectionType {
	return SectionTypeBlock
}

func (s *SectionBlock) GetHorizon(index int, fieldIndexes []int) []string {
	if index >= len(s.data) {
		Exit("[%v] Row index %d is out of length %d", s.sheet, index, len(s.data))
	}
	var contents = make([]string, 0, len(fieldIndexes))
	for i := 0; i < len(fieldIndexes); i++ {
		if fieldIndexes[i] >= len(s.data[index]) {
			Exit("[%v] Row index %d field %d field index %d is out of length %d", s.sheet, index, i, fieldIndexes[i], len(s.data[index]))
		}
		contents = append(contents, s.data[index][fieldIndexes[i]])
	}
	return contents
}
