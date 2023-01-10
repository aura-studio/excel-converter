package converter

import (
	"strings"
)

type SheetRegular struct {
	*SheetBase
}

func NewSheetRegular(excel Excel, name string, data [][]string) *SheetRegular {
	return &SheetRegular{
		SheetBase: NewSheetBase(excel, name, data),
	}
}

func (*SheetRegular) Type() SheetType {
	return SheetTypeRegular
}

func (s *SheetRegular) Read() {
	s.fillSectionMap()
}

func (s *SheetRegular) fillSectionMap() {
	for index, row := range s.data {
		sectionType := s.sectionType(index, row)
		if _, ok := s.sectionMap[sectionType]; !ok {
			s.sectionMap[sectionType] = NewSection(s, sectionType)
		}
		s.sectionMap[sectionType].Append(index, row)
	}
	if _, ok := s.sectionMap[SectionTypeHeader]; !ok {
		Exit("[%v] Sheet has no header", s)
	}
	if _, ok := s.sectionMap[SectionTypeContent]; !ok {
		s.sectionMap[SectionTypeContent] = NewSection(s, SectionTypeContent)
	}
}

func (s *SheetRegular) sectionType(index int, row []string) SectionType {
	switch {
	case index < 2:
		return SectionTypeHeader
	case len(row) == 0 || strings.Contains(row[0], FlagComment):
		return SectionTypeComment
	default:
		return SectionTypeContent
	}
}
