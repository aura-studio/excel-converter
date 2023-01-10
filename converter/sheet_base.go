package converter

import "fmt"

type SheetBase struct {
	excel       Excel
	data        [][]string
	name        string
	sectionMap  map[SectionType]Section
	contentType SectionType
}

func NewSheetBase(excel Excel, name string, data [][]string) *SheetBase {
	return &SheetBase{
		excel:      excel,
		data:       data,
		name:       name,
		sectionMap: map[SectionType]Section{},
	}
}

func (s *SheetBase) String() string {
	return fmt.Sprintf("%v:%s", s.excel, s.name)
}

func (s *SheetBase) Excel() Excel {
	return s.excel
}

func (s *SheetBase) Read() {

}

func (s *SheetBase) FormatHeader(fieldType FieldType) {
	s.GetSection(SectionTypeHeader).Format(0)
	s.GetSection(SectionTypeHeader).(*SectionHeader).ParseFields(fieldType)
}

func (s *SheetBase) FormatContent() {
	s.GetContentSection().Format(s.OriginHeaderSize())
}

func (s *SheetBase) Type() SheetType {
	return ""
}

func (s *SheetBase) Name() string {
	return s.name
}

func (s *SheetBase) FixedName() string {
	return format.ToUpper(s.Name())
}

func (s *SheetBase) IndirectName() string {
	name := s.FixedName()
	if name == "" {
		return FlagDefault
	}
	return name
}

func (s *SheetBase) GetSection(name SectionType) Section {
	if section, ok := s.sectionMap[name]; ok {
		return section
	} else {
		Exit("[%v] Unsupported section name %s", s, name)
	}
	return nil
}

func (s *SheetBase) HeaderSize() int {
	header := s.GetSection(SectionTypeHeader).(*SectionHeader)
	return header.HeaderSize()
}

func (s *SheetBase) OriginHeaderSize() int {
	header := s.GetSection(SectionTypeHeader).(*SectionHeader)
	return header.OriginHeadeSize()
}

func (s *SheetBase) VerticleSize() int {
	return s.GetContentSection().Size()
}

func (s *SheetBase) GetHeaderField(key interface{}) HeaderField {
	header := s.GetSection(SectionTypeHeader).(*SectionHeader)
	return header.GetHeaderField(key)
}

func (s *SheetBase) GetHorizon(index int) []string {
	return s.GetContentSection().GetHorizon(index, s.FieldIndexes())
}

func (s *SheetBase) GetVerticle(index int) []string {
	return s.GetContentSection().GetVertical(index, s.FieldIndexes())
}

func (s *SheetBase) FieldIndexes() []int {
	header := s.GetSection(SectionTypeHeader).(*SectionHeader)
	return header.FieldIndexes()
}

func (s *SheetBase) GetContentSection() Section {
	return s.GetSection(s.contentType)
}

func (s *SheetBase) GetContentKeyIndex(key interface{}) (int, bool) {
	if str, ok := key.(string); ok {
		if str, ok := format.ParseContentKey(str); ok {
			contents := s.GetVerticle(0)
			for index, content := range contents {
				if content == str {
					return index, true
				}
			}
			Exit("[%v] Content key not found %v", s, key)
		}
	}
	return 0, false
}

func (s *SheetBase) GetIndex(key interface{}) int {
	if index, ok := s.GetContentKeyIndex(key); ok {
		return index
	}
	return s.GetContentSection().GetIndex(key)
}

func (s *SheetBase) GetCell(field HeaderField, index int) string {
	return s.GetVerticle(field.Index)[index]
}

func (s *SheetBase) ParseContent(structure StructureType) {
	if s.contentType == "" {
		s.ParseContentType(structure)
		s.UpgradeContent()
		s.FormatContent()
	}
}

func (s *SheetBase) ParseContentType(structureType StructureType) {
	switch structureType {
	case StructureTypeRow, StructureTypeRows, StructureTypeRowMapVal, StructureTypeRowMap:
		s.contentType = SectionTypeRows
	case StructureTypeCol, StructureTypeCols, StructureTypeColMapVal, StructureTypeColMap:
		s.contentType = SectionTypeColumns
	default:
		s.contentType = SectionTypeBlock
	}
}

func (s *SheetBase) UpgradeContent() {
	newSection := NewSection(s, s.contentType)
	section := s.GetSection(SectionTypeContent).(*SectionContent)
	switch newSection.Type() {
	case SectionTypeRows:
		newSection.(*SectionRows).SectionBase = section.SectionBase
	case SectionTypeColumns:
		newSection.(*SectionColumns).SectionBase = section.SectionBase
	case SectionTypeBlock:
		newSection.(*SectionBlock).SectionBase = section.SectionBase
	default:
		Exit("[%v] Unsupported section type %v", s, s.contentType)
	}
	s.sectionMap[newSection.Type()] = newSection
	delete(s.sectionMap, SectionTypeContent)
}

func (s *SheetBase) ForServer() bool {
	return !format.SpecificClient(s.Name())
}

func (s *SheetBase) ForClient() bool {
	return !format.SpecificServer(s.Name())
}
