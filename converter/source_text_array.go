package converter

import (
	"strings"
)

type SourceTextArray struct {
	*SourceBase
	field   *Field
	content string
}

func NewSourceTextArray(ctx *SourceContext, sheet Sheet, field *Field, content string) *SourceTextArray {
	s := &SourceTextArray{
		SourceBase: NewSourceBase(ctx, sheet),
		field:      field,
		content:    content,
	}
	s.ParseSources()
	return s
}

func (s *SourceTextArray) Type() SourceType {
	return SourceTypeTextArray
}

func (s *SourceTextArray) ParseSources() {
	if s.content == "" {
		return
	}
	contents := strings.Split(s.content, FlagComma)
	field := NewField(s.sheet, s.field.Args[0])
	for _, content := range contents {
		content = strings.Trim(content, " ")
		switch field.Structure {
		case StructureTypeBool, StructureTypeInt, StructureTypeFloat, StructureTypeString,
			StructureTypeBigInt, StructureTypeBigFloat, StructureTypeBigRat, StructureTypeTime:
			s.AppendSource(NewSourceTextElement(s.ctx, s.sheet, field, content))
		case StructureTypeStruct, StructureTypeMapVal:
			sheetName := field.Args[0]
			sheet := s.Excel().GetSheet(sheetName)
			if content == "" {
				s.AppendSource(NewSourceNil(s.ctx, sheet))
			} else {
				index := sheet.GetIndex(content)
				s.AppendSource(NewSourceRecordStruct(s.ctx, sheet, index))
			}
		case StructureTypeRow, StructureTypeCol:
			sheetName := field.Args[0]
			sheet := s.Excel().GetSheet(sheetName)
			if content == "" {
				s.AppendSource(NewSourceNil(s.ctx, sheet))
			} else {
				index := sheet.GetIndex(content)
				s.AppendSource(NewSourceRecordSlice(s.ctx, sheet, index))
			}
		case StructureTypeRowMapVal, StructureTypeColMapVal:
			sheetName := field.Args[0]
			sheet := s.Excel().GetSheet(sheetName)
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
					s.AppendSource(NewSourceTable(s.ctx, sheet, field.Structure, nil))
				}
			}
		case StructureTypeArray:
			s.AppendSource(NewSourceTextArray(s.ctx, s.sheet, field, content))
		case StructureTypeTable:
			s.AppendSource(NewSourceTextTable(s.ctx, s.sheet, field, content))
		default:
			Exit("[%v] Unsupported field type %s", s.sheet, field.Structure)
		}
	}
}
