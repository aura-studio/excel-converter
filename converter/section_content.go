package converter

type SectionContent struct {
	*SectionBase
}

func NewSectionContent(sheet Sheet) *SectionContent {
	return &SectionContent{
		SectionBase: NewSectionBase(sheet),
	}
}

// Preprocess 切成方块形
func (s *SectionContent) Preprocess() {

}

func (s *SectionContent) Type() SectionType {
	return SectionTypeContent
}
