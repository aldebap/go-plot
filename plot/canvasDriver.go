////////////////////////////////////////////////////////////////////////////////
//	canvasDriver.go  -  Sep-20-2022  -  aldebap
//
//	Implementation of a graphic driver to generate HTML5 canvas files
////////////////////////////////////////////////////////////////////////////////

package plot

import (
	"bufio"
	"errors"
	"fmt"
)

type Canvas_Driver struct {
	writer       *bufio.Writer
	functionName string
	width        int64
	height       int64
	path         []DriverPoint
	pathColour   RGB_colour
	fontFamily   string
	fontSize     uint8
}

//	NewCanvas_Driver create a new Canvas_Driver
func NewCanvas_Driver(writer *bufio.Writer) GraphicsDriver {
	const (
		PLOT_FUNCTION = "canvas_plot"
		WIDTH         = 600
		HEIGHT        = 400
		FONT_FAMILY   = "Verdana"
		FONT_SIZE     = 10
	)

	return &Canvas_Driver{
		writer:       writer,
		functionName: PLOT_FUNCTION,
		width:        WIDTH,
		height:       HEIGHT,
		path:         nil,
		fontFamily:   FONT_FAMILY,
		fontSize:     FONT_SIZE,
	}
}

//	GetDimensions get the dimensions of the SVG graphic
func (driver *Canvas_Driver) GetDimensions() (width, heigth int64) {
	return driver.width, driver.height
}

//	SetDimensions set the dimensions of the Canvas graphic
func (driver *Canvas_Driver) SetDimensions(width int64, height int64) error {
	driver.width = width
	driver.height = height

	driver.writer.WriteString("function " + driver.functionName + "() {\n" +
		"  let canvas = document.getElementById(\"" + driver.functionName + "\");\n" +
		"  let ctx = canvas.getContext(\"2d\");\n\n")

	driver.Comment("image background")
	driver.writer.WriteString("  ctx.fillStyle = \"#FFFFFF\";\n")
	driver.writer.WriteString("  ctx.fillRect(0, 0, " + fmt.Sprintf("%d", width) + ", " + fmt.Sprintf("%d", height) + ");\n")

	return nil
}

//	GetFont get information about the font
func (driver *Canvas_Driver) GetFont() (fontFamily string, fontSize uint8) {
	return driver.fontFamily, driver.fontSize
}

//	SetFont set information about the font
func (driver *Canvas_Driver) SetFont(fontFamily string, fontSize uint8) {
	driver.fontFamily = fontFamily
	driver.fontSize = fontSize
}

//	Comment write a comment int the the SVG graphic
func (driver *Canvas_Driver) Comment(text string) {
	driver.writer.WriteString("// " + text + "\n")
}

//	Point draws a point in the SVG graphic
func (driver *Canvas_Driver) Point(x, y int64, colour RGB_colour) error {
	driver.writer.WriteString("  ctx.beginPath();\n")
	driver.writer.WriteString("  ctx.strokeStyle = \"#" + colour.Hexa() + "\";\n")
	driver.writer.WriteString("  ctx.moveTo(" + fmt.Sprintf("%d", x) + ", " + fmt.Sprintf("%d", driver.height-y) + ");\n")
	driver.writer.WriteString("  ctx.lineTo(" + fmt.Sprintf("%d", x+1) + ", " + fmt.Sprintf("%d", driver.height-y) + ");\n")
	driver.writer.WriteString("  ctx.stroke();\n")

	return nil
}

//	Begin a path to draw a connection between a set of points
func (driver *Canvas_Driver) BeginPath(colour RGB_colour) error {
	driver.path = make([]DriverPoint, 0)
	driver.pathColour = colour

	return nil
}

//	Add a point to a path
func (driver *Canvas_Driver) PointToPath(x, y int64) error {
	if driver.path == nil {
		return errors.New("cannot add a point to a non initialized path")
	}
	driver.path = append(driver.path, DriverPoint{X: x, Y: y})

	return nil
}

//	End the path
func (driver *Canvas_Driver) EndPath() error {
	if driver.path == nil {
		return errors.New("cannot end a path not initialized")
	}

	driver.writer.WriteString("  ctx.beginPath();\n")
	driver.writer.WriteString("  ctx.strokeStyle = \"#" + driver.pathColour.Hexa() + "\";\n")

	for i, point := range driver.path {
		if i == 0 {
			driver.writer.WriteString("  ctx.moveTo(" + fmt.Sprintf("%d", point.X) + ", " + fmt.Sprintf("%d", driver.height-point.Y) + ");\n")
		} else {
			driver.writer.WriteString("  ctx.lineTo(" + fmt.Sprintf("%d", point.X) + ", " + fmt.Sprintf("%d", driver.height-point.Y) + ");\n")
		}
	}
	driver.writer.WriteString("  ctx.stroke();\n")

	driver.path = nil

	return nil
}

//	Line draws a line between two points in the SVG graphic
func (driver *Canvas_Driver) Line(x1, y1, x2, y2 int64, colour RGB_colour) error {
	driver.writer.WriteString("  ctx.beginPath();\n")
	driver.writer.WriteString("  ctx.strokeStyle = \"#" + colour.Hexa() + "\";\n")
	driver.writer.WriteString("  ctx.moveTo(" + fmt.Sprintf("%d", x1) + ", " + fmt.Sprintf("%d", driver.height-y1) + ");\n")
	driver.writer.WriteString("  ctx.lineTo(" + fmt.Sprintf("%d", x2) + ", " + fmt.Sprintf("%d", driver.height-y2) + ");\n")
	driver.writer.WriteString("  ctx.stroke();\n")

	return nil
}

//	GetTextBox evaluate the width and height of the rectangle required to draw the text string using a given font size
func (driver *Canvas_Driver) GetTextBox(text string) (width, height int64) {

	//	a rough estimation of the rectangle dimentions
	width = int64(0.37 * float64(int64(driver.fontSize)*int64(len(text))))
	height = int64(0.8 * float64(driver.fontSize))

	return width, height
}

//	Text writes a string to the specified point in the SVG graphic
func (driver *Canvas_Driver) Text(x, y, angle int64, text string, colour RGB_colour) error {
	driver.writer.WriteString("  ctx.font = \"" + fmt.Sprintf("%d", driver.fontSize) + "px " + driver.fontFamily + "\";\n")
	driver.writer.WriteString("  ctx.fillStyle = \"#" + colour.Hexa() + "\";\n")
	driver.writer.WriteString("  ctx.fillText(\"" + text + "\", " +
		fmt.Sprintf("%d", x) + ", " + fmt.Sprintf("%d", driver.height-y) + ");\n")

	return nil
}

//	Close finalize the Canvas graphic
func (driver *Canvas_Driver) Close() error {
	driver.writer.WriteString("}\n\n")
	driver.writer.Flush()

	return nil
}
