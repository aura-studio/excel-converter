package main

import (
	"fmt"
	"log"
	"os"
	"rocket-nano/tools/converter"
	"time"
)

func main() {
	if len(os.Args) == 0 {
		log.Println("Requires at least 1 argument")
	}
	fmt.Println("Start...")
	beginTm := time.Now()
	converter.Convert(os.Args[1])
	endTm := time.Now()
	fmt.Printf("Done in %v seconds\n", float64(endTm.UnixNano()-beginTm.UnixNano())/10e8)
}
