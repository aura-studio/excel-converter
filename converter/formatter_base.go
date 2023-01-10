package converter

import "bytes"

type FormatterBase struct {
	*bytes.Buffer
	depth int
}

func NewFormatterBase() *FormatterBase {
	return &FormatterBase{
		Buffer: new(bytes.Buffer),
	}
}

func (f *FormatterBase) FormatIndent() {
	f.WriteString(format.Indent(f.depth))
}

func (f *FormatterBase) IncDepth() {
	f.depth++
}

func (f *FormatterBase) DecDepth() {
	f.depth--
}

func (f *FormatterBase) Depth() int {
	return f.depth
}

func (f *FormatterBase) IsSourceEmpty(source Source) bool {
	for _, source := range source.Sources() {
		if !f.IsSourceEmpty(source) {
			return false
		}
	}
	return source.Empty()
}
