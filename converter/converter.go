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

var converterCreators map[ConverterType]func() Converter

func init() {
	converterCreators = map[ConverterType]func() Converter{
		ConverterTypeGo: func() Converter {
			return NewConverterGo()
		},
		ConverterTypeLua: func() Converter {
			return NewConverterLua()
		},
	}
}

type Config struct {
	Type        string
	ImportPath  string
	ExportPath  string
	ProjectPath string
}

func Run(c Config) {
	if converterCreator, ok := converterCreators[ConverterType(c.Type)]; !ok {
		Exit("[Main] Converter %s is not supported", c.Type)
	} else {
		path.Init(c.ImportPath, c.ExportPath, c.ProjectPath)
		converterCreator().Run()
	}
}
