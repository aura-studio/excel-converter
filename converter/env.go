package converter

type Env struct {
	RenderType RenderType
	DataType   DataType
}

func (e *Env) Init(renderType RenderType, dataType DataType) {
	e.RenderType = renderType
	e.DataType = dataType
}

var env = &Env{}
