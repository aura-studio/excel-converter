package converter

import (
	"strings"
)

type SourceTextTable struct {
	*SourceBase
	field   *Field
	content string
}

func NewSourceTextTable(ctx *SourceContext, sheet Sheet, field *Field, content string) *SourceTextTable {
	s := &SourceTextTable{
		SourceBase: NewSourceBase(ctx, sheet),
		field:      field,
		content:    content,
	}
	s.ParseSources()
	return s
}

func (s *SourceTextTable) Type() SourceType {
	return SourceTypeTextTable
}

func (s *SourceTextTable) ParseSources() {
	if s.content == "" {
		return
	}

	contents := strings.Split(s.content, FlagComma)
	fieldKey := NewField(s.sheet, s.field.Args[0])
	fieldVal := NewField(s.sheet, s.field.Args[1])
	for _, content := range contents {
		content = strings.Trim(content, " ")
		if content == "" {
			continue
		}
		kv := strings.Split(content, FlagColon)
		if len(kv) != 2 {
			Exit("[%v] Source text table %s split error", s.sheet, content)
		}
		switch fieldKey.Structure {
		case StructureTypeBool, StructureTypeInt, StructureTypeFloat, StructureTypeString,
			StructureTypeBigInt, StructureTypeBigFloat, StructureTypeBigRat, StructureTypeTime:
			s.AppendSource(NewSourceTextElement(s.ctx, s.sheet, fieldKey, kv[0]))
		default:
			Exit("[%v] Unsupported field type %s", s.sheet, fieldKey.Structure)
		}
		switch fieldVal.Structure {
		case StructureTypeBool, StructureTypeInt, StructureTypeFloat, StructureTypeString,
			StructureTypeBigInt, StructureTypeBigFloat, StructureTypeBigRat, StructureTypeTime:
			s.AppendSource(NewSourceTextElement(s.ctx, s.sheet, fieldVal, kv[1]))
		case StructureTypeStruct, StructureTypeMapVal:
			sheetName := fieldVal.Args[0]
			sheet := s.Excel().GetSheet(sheetName)
			index := sheet.GetIndex(kv[1])
			s.AppendSource(NewSourceRecordStruct(s.ctx, sheet, index))
		case StructureTypeRow, StructureTypeCol:
			sheetName := fieldVal.Args[0]
			sheet := s.Excel().GetSheet(sheetName)
			index := sheet.GetIndex(kv[1])
			s.AppendSource(NewSourceRecordSlice(s.ctx, sheet, index))
		case StructureTypeRowMapVal, StructureTypeColMapVal:
			sheetName := fieldVal.Args[0]
			sheet := s.Excel().GetSheet(sheetName)
			index := sheet.GetIndex(kv[1])
			s.AppendSource(NewSourceRecordSliceMap(s.ctx, sheet, index))
		case StructureTypeStructs, StructureTypeMap, StructureTypeRows, StructureTypeCols, StructureTypeRowMap, StructureTypeColMap:
			if len(fieldVal.Args) > 0 {
				sheetName := fieldVal.Args[0]
				sheet := s.Excel().GetSheet(sheetName)
				content := kv[1]
				if content == "" {
					s.AppendSource(NewSourceNil(s.ctx, sheet))
				} else {
					s.AppendSource(NewSourceTable(s.ctx, sheet, fieldVal.Structure, s.RangeToIndexes(sheet, content)))
				}
			} else {
				content := kv[1]
				if content == "" {
					s.AppendSource(NewSourceNil(s.ctx, s.sheet))
				} else {
					sheetName := content
					sheet := s.Excel().GetSheet(sheetName)
					s.AppendSource(NewSourceTable(s.ctx, sheet, fieldVal.Structure, nil))
				}
			}
		case StructureTypeArray:
			s.AppendSource(NewSourceTextArray(s.ctx, s.sheet, fieldVal, kv[1]))
		case StructureTypeTable:
			s.AppendSource(NewSourceTextTable(s.ctx, s.sheet, fieldVal, kv[1]))
		default:
			Exit("[%v] Unsupported field type %s", s.sheet, fieldVal.Structure)
		}
	}
}

// table:[int]string
// 1:Tony,2:Sally
