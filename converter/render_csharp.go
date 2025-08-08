package converter

import (
	"path/filepath"
)

type RenderCSharp struct{}

func NewRenderCSharp() *RenderCSharp {
	return &RenderCSharp{}
}

func (r *RenderCSharp) Render(c *Converter) {
	r.FormatStructs(c)
}

func (r *RenderCSharp) FormatStructs(c *Converter) {
	formatter := NewFormatterCSharpStructs(c.identifier)
	formatter.FormatStruct()
	formatter.FormatStructEqual()
	content := formatter.Close()
	filePath := filepath.Join(path.ExportAbsPath(), "Structs.cs")
	c.contentMap[filePath] = content
}
