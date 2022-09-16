////////////////////////////////////////////////////////////////////////////////
//	main.go  -  Sep-4-2022  -  aldebap
//
//	Entry point of a CLI for Go-Plot
////////////////////////////////////////////////////////////////////////////////

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"

	plot "github.com/aldebap/go-plot/plot"
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

	err := generateGraphicFromPlotFile(plotFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[error] fail generating graphic from Go-Plot file: %s\n", err.Error())
		os.Exit(-1)
	}
}

//	generateGraphicFromPlotFile load a Go-Plot file and generate a graphic file from it
func generateGraphicFromPlotFile(plotFileName string) error {

	//	open the Go-Plot file and parse it
	plotFile, err := os.Open(plotFileName)
	if err != nil {
		return errors.New("error opening Go-Plot file: " + err.Error())
	}
	defer plotFile.Close()

	//	load the plot file
	currentPlot, err := plot.LoadPlotFile(bufio.NewReader(plotFile))
	if err != nil {
		return errors.New("fail parsing Go-Plot file: " + err.Error())
	}

	//	create the graphics file for the output
	//	TODO: the output format needs to be decided (or configured!)
	graphicsFile, err := os.Create(plotFileName + ".svg")
	if err != nil {
		return errors.New("error creating graphics file: " + err.Error())
	}
	defer graphicsFile.Close()

	//	create the graphics driver
	driver := plot.NewSVG_Driver(bufio.NewWriter(graphicsFile))
	defer driver.Close()

	//	generate the plot
	err = currentPlot.GeneratePlot(driver)
	if err != nil {
		return errors.New("error generating graphics file: " + err.Error())
	}

	return nil
}
