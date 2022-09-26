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
	writer     *bufio.Writer
	width      int64
	height     int64
	fontFamily string
}

//	create a new SVG_Driver
func NewSVG_Driver(writer *bufio.Writer) GraphicsDriver {
	return &SVG_Driver{
		writer:     writer,
		fontFamily: "Verdana",
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

//	Comment write a comment int the the SVG graphic
func (driver *SVG_Driver) Comment(text string) {
	driver.writer.WriteString("<!-- " + text + "-->\n")
}

func (driver *SVG_Driver) Point(x, y int64, colour RGB_colour) error {
	style := "stroke:rgb(" + fmt.Sprintf("%d", colour.red) +
		"," + fmt.Sprintf("%d", colour.green) +
		"," + fmt.Sprintf("%d", colour.blue) + ");stroke-width:1"

	driver.writer.WriteString("<line x1=\"" + fmt.Sprintf("%d", x) + "\" y1=\"" + fmt.Sprintf("%d", driver.height-y) + "\" " +
		"x2=\"" + fmt.Sprintf("%d", x+1) + "\" y2=\"" + fmt.Sprintf("%d", driver.height-y) + "\" style=\"" + style + "\" />\n")

	return nil
}

//	Line draws a line between two points in the SVG graphic
func (driver *SVG_Driver) Line(x1, y1, x2, y2 int64, colour RGB_colour) error {
	style := "stroke:rgb(" + fmt.Sprintf("%d", colour.red) +
		"," + fmt.Sprintf("%d", colour.green) +
		"," + fmt.Sprintf("%d", colour.blue) + ");stroke-width:1"

	driver.writer.WriteString("<line x1=\"" + fmt.Sprintf("%d", x1) + "\" y1=\"" + fmt.Sprintf("%d", driver.height-y1) + "\" " +
		"x2=\"" + fmt.Sprintf("%d", x2) + "\" y2=\"" + fmt.Sprintf("%d", driver.height-y2) + "\" style=\"" + style + "\" />\n")

	return nil
}

//	GetTextBox evaluate the width and height of the rectangle required to draw the text string using a given font size
func (driver *SVG_Driver) GetTextBox(text string, fontSize int64) (width, height int64) {

	//	a rough estimation of the rectangle dimentions
	width = int64(0.37 * float64(fontSize*int64(len(text))))
	height = int64(0.8 * float64(fontSize))

	return width, height
}

//	Text writes a string to the specified point in the SVG graphic
func (driver *SVG_Driver) Text(x, y, angle int64, text string, fontSize int64, colour RGB_colour) error {
	style := "fill:rgb(" + fmt.Sprintf("%d", colour.red) +
		"," + fmt.Sprintf("%d", colour.green) +
		"," + fmt.Sprintf("%d", colour.blue) + ")"

	if angle == 0 {
		driver.writer.WriteString("<text x=\"" + fmt.Sprintf("%d", x) + "\" y=\"" + fmt.Sprintf("%d", driver.height-y) + "\" " +
			"style=\"" + style + "\" font-family=\"" + driver.fontFamily + "\" font-size=\"" + fmt.Sprintf("%d", fontSize) +
			"\">" + text + "</text>\n")
	} else {
		driver.writer.WriteString("<text x=\"" + fmt.Sprintf("%d", x) + "\" y=\"" + fmt.Sprintf("%d", driver.height-y) + "\" " +
			"transform=\"rotate(" + fmt.Sprintf("%d", angle) + ", " +
			fmt.Sprintf("%d", x) + ", " + fmt.Sprintf("%d", driver.height-y) + ")\" " +
			"style=\"" + style + "\" font-family=\"" + driver.fontFamily + "\" font-size=\"" + fmt.Sprintf("%d", fontSize) +
			"\">" + text + "</text>\n")
	}

	return nil
}

//	Close finalize the SVG graphic
func (driver *SVG_Driver) Close() {
	driver.writer.WriteString("</svg>\n")
	driver.writer.Flush()
}
