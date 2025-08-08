package converter

import "strings"

type HeaderField struct {
	Index int
	Name  string
	Type  string
}

type SectionHeader struct {
	*SectionBase
	fieldMap map[DataType][]int
}

func NewSectionHeader(sheet Sheet) *SectionHeader {
	return &SectionHeader{
		SectionBase: NewSectionBase(sheet),
		fieldMap:    make(map[DataType][]int),
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
	return len(s.fieldMap[DataType(env.DataType)])
}

func (s *SectionHeader) OriginHeadeSize() int {
	return len(s.data[0])
}

func (s *SectionHeader) GetHeaderField(key any) HeaderField {
	switch key := key.(type) {
	case int:
		index := key
		if index >= len(s.data[0]) {
			Exit("[%v] Index is over header size", s.sheet)
		}
		rawIndex := s.fieldMap[env.DataType][index]
		if len(s.data) == 2 {
			return HeaderField{index, format.ToUpper(s.data[0][rawIndex]), s.data[1][rawIndex]}
		} else {
			return HeaderField{index, format.ToUpper(s.data[0][rawIndex]), string(StructureTypeString)}
		}
	case string:
		for index, name := range s.data[0] {
			if key == format.ToUpper(name) {
				rawIndex := s.fieldMap[env.DataType][index]
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

func (s *SectionHeader) ParseFields() {
	for index := 0; index < len(s.data[0]); index++ {
		switch {
		case strings.HasPrefix(s.data[0][index], FlagComment):
			s.fieldMap[DataTypeComment] = append(s.fieldMap[DataTypeComment], index)
		case strings.HasPrefix(s.data[0][index], FlagServer):
			s.fieldMap[DataTypeServer] = append(s.fieldMap[DataTypeServer], index)
		case strings.HasPrefix(s.data[0][index], FlagClient):
			s.fieldMap[DataTypeClient] = append(s.fieldMap[DataTypeClient], index)
		default:
			s.fieldMap[DataTypeServer] = append(s.fieldMap[DataTypeServer], index)
			s.fieldMap[DataTypeClient] = append(s.fieldMap[DataTypeClient], index)
		}
	}
}

func (s *SectionHeader) FieldIndexes() []int {
	return s.fieldMap[env.DataType]
}
