package converter

import (
	"sync/atomic"
)

var sourceIncrease uint32

type SourceBase struct {
	id      uint32
	ctx     *SourceContext
	sheet   Sheet
	sources []Source
}

func NewSourceBase(ctx *SourceContext, sheet Sheet) *SourceBase {
	return &SourceBase{
		id:    atomic.AddUint32(&sourceIncrease, 1),
		ctx:   ctx,
		sheet: sheet,
	}
}

func (*SourceBase) Type() SourceType {
	return ""
}

func (s *SourceBase) Sources() []Source {
	return s.sources
}

func (*SourceBase) Content() string {
	return ""
}

func (s *SourceBase) Sheet() Sheet {
	return s.sheet
}

func (s *SourceBase) Excel() Excel {
	return s.sheet.Excel()
}

func (s *SourceBase) AppendSource(source Source) {
	s.sources = append(s.sources, source)
}

func (s *SourceBase) ParseSources() {
}

func (s *SourceBase) RangeToIndexes(sheet Sheet, content string) []int {
	keys := format.ParseRange(content)
	var indexes = make([]int, 0, len(keys))
	for _, key := range keys {
		indexes = append(indexes, sheet.GetIndex(key))
	}
	return indexes
}

func (s *SourceBase) Empty() bool {
	return true
}
