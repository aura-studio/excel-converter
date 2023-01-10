package converter

import (
	"github.com/spf13/cast"
)

type SourceTextDict struct {
	*SourceBase
	field         *Field
	horizonIndex  int
	verticleIndex int
}

func NewSourceTextDict(ctx *SourceContext, sheet Sheet, field *Field, horizonIndex, verticleIndex int) *SourceTextDict {
	s := &SourceTextDict{
		SourceBase:    NewSourceBase(ctx, sheet),
		field:         field,
		horizonIndex:  horizonIndex,
		verticleIndex: verticleIndex,
	}
	s.ParseVerticleIndex()
	s.ParseSources()
	return s
}

func (s *SourceTextDict) Type() SourceType {
	return SourceTypeTextDict
}

func (s *SourceTextDict) ParseVerticleIndex() {
	if s.field.Index >= 0 {
		s.verticleIndex = s.field.Index + 1
	}
}

func (s *SourceTextDict) ParseSources() {
	verticleIndex := s.verticleIndex
	contents := s.sheet.GetHorizon(s.horizonIndex)
	for _, typArg := range s.field.Args {
		if size, err := cast.ToIntE(typArg); err == nil {
			for index := 0; index < size; {
				field := NewField(s.sheet, s.sheet.GetHeaderField(verticleIndex))
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
		} else {
			field := NewField(s.sheet, typArg)
			switch field.Structure {
			case StructureTypeList:
				s.AppendSource(NewSourceTextList(s.ctx, s.sheet, field, s.horizonIndex, verticleIndex))
			case StructureTypeDict:
				s.AppendSource(NewSourceTextDict(s.ctx, s.sheet, field, s.horizonIndex, verticleIndex))
			default:
				Exit("[%v] Unsupported field type %s", s.sheet, field.Structure)
			}
			verticleIndex += field.Size
		}
	}
}
