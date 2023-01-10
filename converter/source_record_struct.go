package converter

type SourceRecordStruct struct {
	*SourceBase
	index int
}

func NewSourceRecordStruct(ctx *SourceContext, sheet Sheet, index int) *SourceRecordStruct {
	s := &SourceRecordStruct{
		SourceBase: NewSourceBase(ctx, sheet),
		index:      index,
	}
	s.ParseSources()
	return s
}

func (s *SourceRecordStruct) Type() SourceType {
	return SourceTypeRecordStruct
}

func (s *SourceRecordStruct) ParseSources() {
	contents := s.sheet.GetHorizon(s.index)
	for index := 0; index < len(contents); {
		headerField := s.sheet.GetHeaderField(index)
		field := NewField(s.sheet, headerField)
		content := contents[index]
		switch field.Structure {
		case StructureTypeBool, StructureTypeInt, StructureTypeFloat, StructureTypeString,
			StructureTypeBigInt, StructureTypeBigFloat, StructureTypeBigRat, StructureTypeTime:
			s.AppendSource(NewSourceTextElement(s.ctx, s.sheet, field, s.sheet.GetCell(headerField, s.index)))
		case StructureTypeStruct, StructureTypeMapVal:
			sheetName := field.Args[0]
			sheet := s.Excel().GetSheet(sheetName)
			sheet.ParseContent(field.Structure)
			if content == "" {
				s.AppendSource(NewSourceNil(s.ctx, sheet))
			} else {
				s.AppendSource(NewSourceRecordStruct(s.ctx, sheet, sheet.GetIndex(content)))
			}
		case StructureTypeRow, StructureTypeCol:
			sheetName := field.Args[0]
			sheet := s.Excel().GetSheet(sheetName)
			sheet.ParseContent(field.Structure)
			if content == "" {
				s.AppendSource(NewSourceNil(s.ctx, sheet))
			} else {
				s.AppendSource(NewSourceRecordSlice(s.ctx, sheet, sheet.GetIndex(content)))
			}
		case StructureTypeRowMapVal, StructureTypeColMapVal:
			sheetName := field.Args[0]
			sheet := s.Excel().GetSheet(sheetName)
			sheet.ParseContent(field.Structure)
			if content == "" {
				s.AppendSource(NewSourceNil(s.ctx, sheet))
			} else {
				s.AppendSource(NewSourceRecordSliceMap(s.ctx, sheet, sheet.GetIndex(content)))
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
