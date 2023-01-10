package converter

type SectionComment struct {
	*SectionBase
}

func NewSectionComment(sheet Sheet) *SectionComment {
	return &SectionComment{
		SectionBase: NewSectionBase(sheet),
	}
}

func (s *SectionComment) Type() SectionType {
	return SectionTypeComment
}
