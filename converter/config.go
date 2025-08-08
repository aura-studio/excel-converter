package converter

import (
	"flag"
	"fmt"
)

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
