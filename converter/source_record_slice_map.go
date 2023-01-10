package converter

type SourceRecordSliceMap struct {
	*SourceBase
	index int
}

func NewSourceRecordSliceMap(ctx *SourceContext, sheet Sheet, index int) *SourceRecordSliceMap {
	s := &SourceRecordSliceMap{
		SourceBase: NewSourceBase(ctx, sheet),
		index:      index,
	}
	s.ParseSources()
	return s
}

func (s *SourceRecordSliceMap) Type() SourceType {
	return SourceTypeRecordSliceMap
}

func (s *SourceRecordSliceMap) ParseSources() {
	headerField := s.sheet.GetHeaderField(1)
	field := NewField(s.sheet, headerField)
	contents := s.sheet.GetHorizon(s.index)[1:]
	for index := 0; index < len(contents); {
		content := contents[index]
		switch field.Structure {
		case StructureTypeBool, StructureTypeInt, StructureTypeFloat, StructureTypeString,
			StructureTypeBigInt, StructureTypeBigFloat, StructureTypeBigRat, StructureTypeTime:
			s.AppendSource(NewSourceTextElement(s.ctx, s.sheet, field, content))
		case StructureTypeStruct, StructureTypeMapVal:
			sheetName := field.Args[0]
			sheet := s.Excel().GetSheet(sheetName)
			sheet.ParseContent(field.Structure)
			if content == "" {
				s.AppendSource(NewSourceNil(s.ctx, sheet))
			} else {
				index := sheet.GetIndex(content)
				s.AppendSource(NewSourceRecordStruct(s.ctx, sheet, index))
			}
		case StructureTypeRow, StructureTypeCol:
			sheetName := field.Args[0]
			sheet := s.Excel().GetSheet(sheetName)
			sheet.ParseContent(field.Structure)
			index := sheet.GetIndex(content)
			s.AppendSource(NewSourceRecordSlice(s.ctx, sheet, index))
		case StructureTypeRowMapVal, StructureTypeColMapVal:
			sheetName := field.Args[0]
			sheet := s.Excel().GetSheet(sheetName)
			sheet.ParseContent(field.Structure)
			if content == "" {
				s.AppendSource(NewSourceNil(s.ctx, sheet))
			} else {
				index := sheet.GetIndex(content)
				s.AppendSource(NewSourceRecordSliceMap(s.ctx, sheet, index))
			}
		case StructureTypeStructs, StructureTypeMap, StructureTypeRows, StructureTypeCols, StructureTypeRowMap, StructureTypeColMap:
			if len(field.Args) > 0 {
				sheetName := field.Args[0]
				sheet := s.Excel().GetSheet(sheetName)
				sheet.ParseContent(field.Structure)
				if content == "" {
					s.AppendSource(NewSourceNil(s.ctx, sheet))
				} else {
					s.AppendSource(NewSourceTable(s.ctx, sheet, field.Structure, s.RangeToIndexes(sheet, content)))
				}
			} else {
				if content == "" {
					s.AppendSource(NewSourceNil(s.ctx, s.sheet))
				} else {
					sheetName := content
					sheet := s.Excel().GetSheet(sheetName)
					sheet.ParseContent(field.Structure)
					s.AppendSource(NewSourceTable(s.ctx, sheet, field.Structure, nil))
				}
			}
		case StructureTypeArray:
			s.AppendSource(NewSourceTextArray(s.ctx, s.sheet, field, content))
		case StructureTypeTable:
			s.AppendSource(NewSourceTextTable(s.ctx, s.sheet, field, content))
		case StructureTypeList:
			s.AppendSource(NewSourceTextList(s.ctx, s.sheet, field, s.index, index))
		case StructureTypeDict:
			s.AppendSource(NewSourceTextDict(s.ctx, s.sheet, field, s.index, index))
		default:
			Exit("[%v] Unsupported field type %s", s.sheet, field.Structure)
		}
		index += field.Size
	}
}
