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
	width  int64
	height int64
}

//	create a new SVG_Driver
func NewSVG_Driver(writer *bufio.Writer) GraphicsDriver {
	return &SVG_Driver{
		writer: writer,
	}
}

//	SetDimensions set the dimensions of the SVG graphic
func (driver *SVG_Driver) SetDimensions(width int64, height int64) error {
	driver.width = width
	driver.height = height

	driver.writer.WriteString("<svg width=\"" + fmt.Sprintf("%d", width) + "\" height=\"" + fmt.Sprintf("%d", height) + "\" " +
		"xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\">\n")

	return nil
}

func (driver *SVG_Driver) Point(x, y int64) error {
	return nil
}

func (driver *SVG_Driver) Line(x1, y1, x2, y2 int64) error {
	const style = "stroke:rgb(0,0,0);stroke-width:1"

	driver.writer.WriteString("<line x1=\"" + fmt.Sprintf("%d", driver.height-x1) + "\" y1=\"" + fmt.Sprintf("%d", y1) + "\" " +
		"x2=\"" + fmt.Sprintf("%d", driver.height-x2) + "\" y2=\"" + fmt.Sprintf("%d", y2) + "\" style=\"" + style + "\" />\n")

	return nil
}

//	Close finalize the SVG graphic
func (driver *SVG_Driver) Close() {
	driver.writer.WriteString("</svg>\n")
	driver.writer.Flush()
}
