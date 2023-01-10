package converter

import "strings"

type FieldType string

const (
	FieldTypeServer  FieldType = "server"
	FieldTypeClient  FieldType = "client"
	FieldTypeComment FieldType = "comment"
)

type HeaderField struct {
	Index int
	Name  string
	Type  string
}

type SectionHeader struct {
	*SectionBase
	fieldMap         map[FieldType][]int
	currentFieldType FieldType
}

func NewSectionHeader(sheet Sheet) *SectionHeader {
	return &SectionHeader{
		SectionBase: NewSectionBase(sheet),
		fieldMap:    make(map[FieldType][]int),
	}
}

func (s *SectionHeader) Type() SectionType {
	return SectionTypeHeader
}

func (s *SectionHeader) Format(int) {
	s.TrimCells()
	s.TrimColumns()
	s.CutRows()
	s.CheckSize()
}

func (s *SectionHeader) CheckSize() {
	switch size := len(s.data); size {
	case 0:
		Exit("[%v] Header is empty", s.sheet)
	case 1:
		return
	case 2:
		if len(s.data[0]) != len(s.data[1]) {
			Exit("[%v] Header name size does not equal to type size", s.sheet)
		}
	case 3:
		Exit("[%v] Header is over 2 rows", s.sheet)
	}
}

func (s *SectionHeader) HeaderSize() int {
	return len(s.fieldMap[s.currentFieldType])
}

func (s *SectionHeader) OriginHeadeSize() int {
	return len(s.data[0])
}

func (s *SectionHeader) GetHeaderField(key interface{}) HeaderField {
	switch key := key.(type) {
	case int:
		index := key
		if index >= len(s.data[0]) {
			Exit("[%v] Index is over header size", s.sheet)
		}
		rawIndex := s.fieldMap[s.currentFieldType][index]
		if len(s.data) == 2 {
			return HeaderField{index, format.ToUpper(s.data[0][rawIndex]), s.data[1][rawIndex]}
		} else {
			return HeaderField{index, format.ToUpper(s.data[0][rawIndex]), string(StructureTypeString)}
		}
	case string:
		for index, name := range s.data[0] {
			if key == format.ToUpper(name) {
				rawIndex := s.fieldMap[s.currentFieldType][index]
				if len(s.data) == 2 {
					return HeaderField{index, format.ToUpper(s.data[0][rawIndex]), s.data[1][rawIndex]}
				} else {
					return HeaderField{index, format.ToUpper(s.data[0][rawIndex]), string(StructureTypeString)}
				}
			}
		}
	}

	Exit("[%v] Unsupported field key %v", s.sheet, key)
	return HeaderField{}
}

func (s *SectionHeader) fmtFieldName(name string) string {
	return format.ToUpper(name)
}

func (s *SectionHeader) ParseFields(fieldType FieldType) {
	for index := 0; index < len(s.data[0]); index++ {
		switch {
		case strings.HasPrefix(s.data[0][index], FlagComment):
			s.fieldMap[FieldTypeComment] = append(s.fieldMap[FieldTypeComment], index)
		case strings.HasPrefix(s.data[0][index], FlagServer):
			s.fieldMap[FieldTypeServer] = append(s.fieldMap[FieldTypeServer], index)
		case strings.HasPrefix(s.data[0][index], FlagClient):
			s.fieldMap[FieldTypeClient] = append(s.fieldMap[FieldTypeClient], index)
		default:
			s.fieldMap[FieldTypeServer] = append(s.fieldMap[FieldTypeServer], index)
			s.fieldMap[FieldTypeClient] = append(s.fieldMap[FieldTypeClient], index)
		}
	}
	s.currentFieldType = fieldType
}

func (s *SectionHeader) FieldIndexes() []int {
	return s.fieldMap[s.currentFieldType]
}
