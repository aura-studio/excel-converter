package main

import (
	"fmt"
	"log"
	"os"
	"time"

	. "github.com/aura-studio/excel-converter/converter"
)

func main() {
	if len(os.Args) < 3 {
		log.Println("Requires at least 3 argument")
	}
	fmt.Println("Start...")
	beginTm := time.Now()

	defer func() {
		endTm := time.Now()
		fmt.Printf("Done in %v seconds\n", float64(endTm.UnixNano()-beginTm.UnixNano())/10e8)
	}()

	c := Config{
		Type:       os.Args[1],
		ImportPath: os.Args[2],
		ExportPath: os.Args[3],
	}

	if len(os.Args) > 4 {
		c.ProjectPath = os.Args[4]
	}

	Run(c)
}
