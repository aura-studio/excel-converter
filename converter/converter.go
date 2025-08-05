package converter

import (
	"flag"
	"fmt"
)

type ConverterType string

const (
	ConverterTypeGo    ConverterType = "go"
	ConverterTypeLua   ConverterType = "lua"
	ConverterTypeJson  ConverterType = "json"
	ConverterTypeProto ConverterType = "proto"
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
		ConverterTypeJson: func() Converter {
			return NewConverterJson()
		},
	}
}

type Config struct {
	Type        string
	ImportPath  string
	ExportPath  string
	ProjectPath string
}

func (c *Config) Parse() {
	if len(flag.Args()) < 5 {
		c.Type = flag.Arg(0)
		c.ImportPath = flag.Arg(1)
		c.ExportPath = flag.Arg(2)
		if c.Type == "go" {
			c.ProjectPath = flag.Arg(3)
		}

		return
	}

	flag.StringVar(&c.Type, "type", "", "type of converter")
	flag.StringVar(&c.ImportPath, "import", "", "import path of excel files")
	flag.StringVar(&c.ExportPath, "export", "", "export path of generated files")
	flag.StringVar(&c.ProjectPath, "project", "", "project path of generated files")
	flag.Parse()
}

func (c *Config) Assert() {
	if c.Type == "" {
		fmt.Println("Type is required")
		return
	}

	if c.ImportPath == "" {
		fmt.Println("ImportPath is required")
		return
	}

	if c.ExportPath == "" {
		fmt.Println("ExportPath is required")
		return
	}

	if c.Type == "go" && c.ProjectPath == "" {
		fmt.Println("ProjectPath is required")
		return
	}
}

func Run(c Config) {
	if converterCreator, ok := converterCreators[ConverterType(c.Type)]; !ok {
		Exit("[Main] Converter %s is not supported", c.Type)
	} else {
		path.Init(c.ImportPath, c.ExportPath, c.ProjectPath)
		converterCreator().Run()
	}
}
