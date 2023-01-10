package converter

import "github.com/urfave/cli"

func Convert(arg string) {
	converterType := ConverterType(arg)
	if converterCreator, ok := converterCreators[converterType]; !ok {
		Exit("[Main] Converter %s is not supported", converterType)
	} else {
		converterCreator(NewPathExternal()).Run()
	}
}

func Run(context *cli.Context) error {
	args := context.Args()
	Convert(args[0])
	return nil
}
