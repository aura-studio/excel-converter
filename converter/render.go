package converter

type Render interface {
	Render(*Converter)
}

type RenderType = string

const (
	RenderTypeGo     RenderType = "go"
	RenderTypeLua    RenderType = "lua"
	RenderTypeJson   RenderType = "json"
	RenderTypeCSharp RenderType = "csharp"
)

var renderMap = map[RenderType]Render{
	RenderTypeGo:     NewRenderGo(),
	RenderTypeLua:    NewRenderLua(),
	RenderTypeJson:   NewRenderJson(),
	RenderTypeCSharp: NewRenderCSharp(),
}
