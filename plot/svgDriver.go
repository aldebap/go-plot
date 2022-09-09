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
func NewSVG_Driver(writer *bufio.Writer) GraphicsDriver {
	return &SVG_Driver{
		writer: writer,
	}
}

//	SetDimensions set the dimensions of the SVG graphic
func (driver *SVG_Driver) SetDimensions(width int64, height int64) error {
	driver.writer.WriteString("<svg width=\"" + fmt.Sprintf("%d", width) + "\" height=\"" + fmt.Sprintf("%d", height) + "\">\n")

	return nil
}

func (driver *SVG_Driver) Point(x, y int64) error {
	return nil
}

func (driver *SVG_Driver) Line(x1, y1, x2, y2 int64) error {
	return nil
}

//	Close finalize the SVG graphic
func (driver *SVG_Driver) Close() {
	driver.writer.WriteString("</svg>\n")
	driver.writer.Flush()
}
