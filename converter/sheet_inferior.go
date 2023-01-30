package converter

import (
	"strings"
)

type SheetInferior struct {
	*SheetBase
}

func NewSheetInferior(excel Excel, name string, data [][]string) *SheetInferior {
	return &SheetInferior{
		SheetBase: NewSheetBase(excel, name, data),
	}
}

func (*SheetInferior) Type() SheetType {
	return SheetTypeInferior
}

func (s *SheetInferior) Read() {
	s.fillSectionMap()
}

func (s *SheetInferior) fillSectionMap() {
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

func (s *SheetInferior) sectionType(index int, row []string) SectionType {
	switch {
	case index < 2:
		return SectionTypeHeader
	case len(row) == 0 || strings.Contains(row[0], FlagComment):
		return SectionTypeComment
	default:
		return SectionTypeContent
	}
}
