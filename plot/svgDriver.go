////////////////////////////////////////////////////////////////////////////////
//	svgDriver.go  -  Sep-8-2022  -  aldebap
//
//	Implementation of a graphic driver to generate SVG files
////////////////////////////////////////////////////////////////////////////////

package plot

import (
	"bufio"
	"fmt"
)

type SVG_Driver struct {
	writer *bufio.Writer
}

//	create a new SVG_Driver
func NewSVG_Driver(writer *bufio.Writer, width int64, height int64) GraphicsDriver {

	writer.WriteString("<svg width=\"" + fmt.Sprintf("%d", width) + "\" height=\"" + fmt.Sprintf("%d", height) + "\">\n")

	return &SVG_Driver{
		writer: writer,
	}
}

func (driver *SVG_Driver) Point(x, y int64) error {
	return nil
}

func (driver *SVG_Driver) Line(x1, y1, x2, y2 int64) error {
	return nil
}
