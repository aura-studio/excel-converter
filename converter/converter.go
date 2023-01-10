package converter

type ConverterType string

const (
	ConverterTypeGo  ConverterType = "go"
	ConverterTypeLua ConverterType = "lua"
)

type Converter interface {
	String() string
	Run()
}

var converterCreators map[ConverterType]func(Path) Converter

func init() {
	converterCreators = map[ConverterType]func(Path) Converter{
		ConverterTypeGo: func(path Path) Converter {
			return NewConverterGo(path)
		},
		ConverterTypeLua: func(path Path) Converter {
			return NewConverterLua(path)
		},
	}
}
