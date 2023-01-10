package converter

import "github.com/mohae/deepcopy"

type SourceTable struct {
	*SourceBase
	indexes   []int
	structure StructureType
}

func NewSourceTable(ctx *SourceContext, sheet Sheet, structure StructureType, indexes []int) *SourceTable {
	s := &SourceTable{
		SourceBase: NewSourceBase(ctx, sheet),
		indexes:    indexes,
		structure:  structure,
	}
	s.ParseSources()
	return s
}

func (s *SourceTable) Type() SourceType {
	return SourceTypeTable
}

func (s *SourceTable) ParseSources() {
	s.sheet.ParseContent(s.structure)
	switch s.structure {
	case StructureTypeStructs:
		if len(s.indexes) == 0 {
			for index := 0; index < s.sheet.VerticleSize(); index++ {
				s.AppendSource(NewSourceRecordStruct(s.ctx, s.sheet, index))
			}
		} else {
			for _, index := range s.indexes {
				s.AppendSource(NewSourceRecordStruct(s.ctx, s.sheet, index))
			}
		}
	case StructureTypeMap:
		if len(s.indexes) == 0 {
			for index := 0; index < s.sheet.VerticleSize(); index++ {
				source := NewSourceRecordStruct(s.ctx, s.sheet, index)
				ctx := deepcopy.Copy(s.ctx).(*SourceContext)
				ctx.key = ""
				if s.ctx.key == "" {
					s.AppendSource(source.Sources()[0])
				} else {
					headerField := s.sheet.GetHeaderField(s.ctx.key)
					s.AppendSource(source.Sources()[headerField.Index])
				}

				s.AppendSource(source)
			}
		} else {
			for _, index := range s.indexes {
				source := NewSourceRecordStruct(s.ctx, s.sheet, index)
				s.AppendSource(source.Sources()[0])
				s.AppendSource(source)
			}
		}
	case StructureTypeRowMap, StructureTypeColMap:
		if len(s.indexes) == 0 {
			for index := 0; index < s.sheet.VerticleSize(); index++ {
				field := NewField(s.sheet, s.sheet.GetHeaderField(0))
				content := s.sheet.GetHorizon(index)[0]
				s.AppendSource(NewSourceTextElement(s.ctx, s.sheet, field, content))
				s.AppendSource(NewSourceRecordSliceMap(s.ctx, s.sheet, index))
			}
		} else {
			for _, index := range s.indexes {
				field := NewField(s.sheet, s.sheet.GetHeaderField(0))
				content := s.sheet.GetHorizon(index)[0]
				s.AppendSource(NewSourceTextElement(s.ctx, s.sheet, field, content))
				s.AppendSource(NewSourceRecordSliceMap(s.ctx, s.sheet, index))
			}
		}
	case StructureTypeRows, StructureTypeCols:
		if len(s.indexes) == 0 {
			for index := 0; index < s.sheet.VerticleSize(); index++ {
				s.AppendSource(NewSourceRecordSlice(s.ctx, s.sheet, index))
			}
		} else {
			for _, index := range s.indexes {
				s.AppendSource(NewSourceRecordSlice(s.ctx, s.sheet, index))
			}
		}

	default:
		Exit("[%v] Unsupported structure type %s", s.sheet, s.structure)
	}
}
