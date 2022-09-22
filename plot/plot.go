////////////////////////////////////////////////////////////////////////////////
//	plot.go  -  Sep-4-2022  -  aldebap
//
//	Interface for a general Go-Plot
////////////////////////////////////////////////////////////////////////////////

package plot

import "bufio"

//	terminals to generate a plot
const (
	TERMINAL_CANVAS = 1
	TERMINAL_PNG    = 2
	TERMINAL_SVG    = 3
)

type Plot interface {
	GetOutputFileName() string
	GeneratePlot(plotWriter *bufio.Writer) error
}
