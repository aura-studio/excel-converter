package converter

import (
	"strings"
)

type SheetSettings struct {
	*SheetBase
}

func NewSheetSettings(excel Excel, name string, data [][]string) *SheetSettings {
	return &SheetSettings{
		SheetBase: NewSheetBase(excel, name, data),
	}
}

func (*SheetSettings) Type() SheetType {
	return SheetTypeSettings
}

func (s *SheetSettings) Read() {
	s.fillSectionMap()
}

func (s *SheetSettings) sectionType(index int, row []string) SectionType {
	switch {
	case index < 1:
		return SectionTypeHeader
	case len(row) == 0 || strings.Contains(row[0], FlagComment):
		return SectionTypeComment
	default:
		return SectionTypeContent
	}
}

func (s *SheetSettings) fillSectionMap() {
	for index, row := range s.data {
		sectionType := s.sectionType(index, row)
		if _, ok := s.sectionMap[sectionType]; !ok {
			s.sectionMap[sectionType] = NewSection(s, sectionType)
		}
		s.sectionMap[sectionType].Append(index, row)
	}
}
