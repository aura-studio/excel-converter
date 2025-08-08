package converter

import (
	"path/filepath"
)

type RenderCSharp struct{}

func NewRenderCSharp() *RenderCSharp {
	return &RenderCSharp{}
}

func (r *RenderCSharp) Render() {
	r.FormatStructs()
}

func (r *RenderCSharp) FormatStructs() {
	formatter := NewFormatterCSharpStructs(c.identifier)
	formatter.FormatStruct()
	formatter.FormatStructEqual()
	content := formatter.Close()
	filePath := filepath.Join(path.ExportAbsPath(), "Structs.cs")
	c.contentMap[filePath] = content
}

// func (r *RenderCSharp) FormatJsonData(c *Con
