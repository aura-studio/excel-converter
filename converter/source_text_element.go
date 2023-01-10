package converter

type SourceTextElement struct {
	*SourceBase
	field   *Field
	content string
}

func NewSourceTextElement(ctx *SourceContext, sheet Sheet, field *Field, content string) *SourceTextElement {
	s := &SourceTextElement{
		SourceBase: NewSourceBase(ctx, sheet),
		field:      field,
		content:    content,
	}
	s.ParseSources()
	return s
}

func (s *SourceTextElement) Empty() bool {
	return s.content == ""
}

func (s *SourceTextElement) Type() SourceType {
	return SourceTypeTextElement
}

func (s *SourceTextElement) Content() string {
	return s.content
}
