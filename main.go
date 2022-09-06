////////////////////////////////////////////////////////////////////////////////
//	main.go  -  Sep-4-2022  -  aldebap
//
//	Entry point of a CLI for Go-Plot
////////////////////////////////////////////////////////////////////////////////

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/aldebap/go-plot/plot"
)

const (
	versionInfo string = "Go-Plot 0.1"
)

//	main entry point for Go-Plot CLI
func main() {
	var (
		version bool
	)

	//	CLI arguments
	flag.BoolVar(&version, "version", false, "show Go-Plot version")

	flag.Parse()

	//	version option
	if version {
		fmt.Printf("%s\n", versionInfo)
		return
	}

	//	get the Go-Plot file name
	plotFileName := flag.Arg(0)
	if len(plotFileName) == 0 {
		fmt.Fprintf(os.Stderr, "[error] missing Go-Plot file name\n")
		os.Exit(-1)
	}

	//	open the Go-Plot file and parse it
	plotFile, err := os.Open(plotFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[error] fail attempting to open Go-Plot file: %s\n", err.Error())
		os.Exit(-1)
	}
	defer plotFile.Close()

	//plotInfo, err := plot.LoadPlotFile(bufio.NewReader(plotFile))
	_, err = plot.LoadPlotFile(bufio.NewReader(plotFile))
	if err != nil {
		fmt.Fprintf(os.Stderr, "[error] fail attempting to parse Go-Plot file: %s\n", err.Error())
		os.Exit(-1)
	}
}
