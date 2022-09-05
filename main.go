////////////////////////////////////////////////////////////////////////////////
//	main.go  -  Sep-4-2022  -  aldebap
//
//	Entry point of a CLI for Go-Plot
////////////////////////////////////////////////////////////////////////////////

package main

import (
	"flag"
	"fmt"
)

const (
	versionNumber string = "0.1"
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
		fmt.Printf("Go-Plot %s\n", versionNumber)
		return
	}
}
