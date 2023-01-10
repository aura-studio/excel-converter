package converter

import (
	"github.com/spf13/cast"
)

type SourceTextList struct {
	*SourceBase
	field         *Field
	horizonIndex  int
	verticleIndex int
}

func NewSourceTextList(ctx *SourceContext, sheet Sheet, field *Field, horizonIndex, verticleIndex int) *SourceTextList {
	s := &SourceTextList{
		SourceBase:    NewSourceBase(ctx, sheet),
		field:         field,
		horizonIndex:  horizonIndex,
		verticleIndex: verticleIndex,
	}
	s.ParseVerticleIndex()
	s.ParseSources()
	return s
}

func (s *SourceTextList) Type() SourceType {
	return SourceTypeTextList
}

func (s *SourceTextList) ParseVerticleIndex() {
	if s.field.Index >= 0 {
		s.verticleIndex = s.field.Index + 1
	}
}

func (s *SourceTextList) ParseSources() {
	verticleIndex := s.verticleIndex
	field := NewField(s.sheet, s.field.Args[0])
	contents := s.sheet.GetHorizon(s.horizonIndex)
	size := cast.ToInt(s.field.Args[1])
	for index := 0; index < size; {
		content := contents[verticleIndex]
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
		case StructureTypeList:
			s.AppendSource(NewSourceTextList(s.ctx, s.sheet, field, s.horizonIndex, verticleIndex))
		case StructureTypeDict:
			s.AppendSource(NewSourceTextDict(s.ctx, s.sheet, field, s.horizonIndex, verticleIndex))
		default:
			Exit("[%v] Unsupported field type %s", s.sheet, field.Structure)
		}
		index++
		verticleIndex += field.Size
	}
}
