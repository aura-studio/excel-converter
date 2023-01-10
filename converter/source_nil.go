package converter

type SourceNil struct {
	*SourceBase
}

func NewSourceNil(ctx *SourceContext, sheet Sheet) *SourceNil {
	s := &SourceNil{
		SourceBase: NewSourceBase(ctx, sheet),
	}
	return s
}

func (s *SourceNil) Type() SourceType {
	return SourceTypeNil
}
