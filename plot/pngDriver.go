////////////////////////////////////////////////////////////////////////////////
//	pngDriver.go  -  Nov-24-2022  -  aldebap
//
//	Implementation of a graphic driver to generate PNG files
////////////////////////////////////////////////////////////////////////////////

package plot

import (
	"bufio"
	"errors"
	"image"
	"image/color"
	"image/png"
)

type PNG_Driver struct {
	writer     *bufio.Writer
	width      int64
	height     int64
	image      *image.RGBA
	path       []DriverPoint
	pathColour RGB_colour
	fontFamily string
	fontSize   uint8
}

//	create a new PNG_Driver
func NewPNG_Driver(writer *bufio.Writer) GraphicsDriver {
	const (
		WIDTH       = 640
		HEIGHT      = 480
		FONT_FAMILY = "Verdana"
		FONT_SIZE   = 10
	)

	return &PNG_Driver{
		writer:     writer,
		width:      WIDTH,
		height:     HEIGHT,
		fontFamily: FONT_FAMILY,
		fontSize:   FONT_SIZE,
	}
}

//	GetDimensions get the dimensions of the PNG graphic
func (driver *PNG_Driver) GetDimensions() (width, heigth int64) {
	return driver.width, driver.height
}

//	SetDimensions set the dimensions of the PNG graphic
func (driver *PNG_Driver) SetDimensions(width int64, height int64) error {
	driver.width = width
	driver.height = height

	driver.image = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(width), int(height)}})

	return nil
}

//	GetFont get information about the font
func (driver *PNG_Driver) GetFont() (fontFamily string, fontSize uint8) {
	return driver.fontFamily, driver.fontSize
}

//	SetFont set information about the font
func (driver *PNG_Driver) SetFont(fontFamily string, fontSize uint8) {
	driver.fontFamily = fontFamily
	driver.fontSize = fontSize
}

//	Comment write a comment int the the PNG graphic
func (driver *PNG_Driver) Comment(text string) {
	//	cannot add comments to PNG
}

//	Point draws a point in the PNG graphic
func (driver *PNG_Driver) Point(x, y int64, colour RGB_colour) error {
	pointColour := color.RGBA{colour.red, colour.green, colour.blue, 255}
	driver.image.Set(int(x), int(y), pointColour)

	return nil
}

//	Begin a path to draw a connection between a set of points
func (driver *PNG_Driver) BeginPath(colour RGB_colour) error {
	driver.path = make([]DriverPoint, 0)
	driver.pathColour = colour

	return nil
}

//	Add a point to a path
func (driver *PNG_Driver) PointToPath(x, y int64) error {
	if driver.path == nil {
		return errors.New("cannot add a point to a non initialized path")
	}
	driver.path = append(driver.path, DriverPoint{X: x, Y: y})

	return nil
}

//	End the path
func (driver *PNG_Driver) EndPath() error {
	if driver.path == nil {
		return errors.New("cannot end a path not initialized")
	}

	/*
		style := "stroke:rgb(" + fmt.Sprintf("%d", driver.pathColour.red) +
			"," + fmt.Sprintf("%d", driver.pathColour.green) +
			"," + fmt.Sprintf("%d", driver.pathColour.blue) + ");stroke-width:1"

		for i, point := range driver.path {
			if i == 0 {
				driver.writer.WriteString("<path d=\"M" + fmt.Sprintf("%d", point.X) + " " + fmt.Sprintf("%d", driver.height-point.Y))
			} else {
				driver.writer.WriteString(" L" + fmt.Sprintf("%d", point.X) + " " + fmt.Sprintf("%d", driver.height-point.Y))
			}
		}
		driver.writer.WriteString("\" style=\"" + style + "\" fill=\"none\" />\n")
	*/
	driver.path = nil

	return nil
}

//	Line draws a line between two points in the PNG graphic
func (driver *PNG_Driver) Line(x1, y1, x2, y2 int64, colour RGB_colour) error {
	lineColour := color.RGBA{colour.red, colour.green, colour.blue, 255}

	y := float64(y1)
	delta := float64(x2-x1) / float64(y2-y1)
	for x := x1; x <= x2; x++ {
		driver.image.Set(int(x), int(y), lineColour)
		y += delta
	}

	return nil
}

//	GetTextBox evaluate the width and height of the rectangle required to draw the text string using a given font size
func (driver *PNG_Driver) GetTextBox(text string) (width, height int64) {

	//	a rough estimation of the rectangle dimentions
	width = int64(0.37 * float64(int64(driver.fontSize)*int64(len(text))))
	height = int64(0.8 * float64(driver.fontSize))

	return width, height
}

//	Text writes a string to the specified point in the PNG graphic
func (driver *PNG_Driver) Text(x, y, angle int64, text string, colour RGB_colour) error {
	/*
		style := "fill:rgb(" + fmt.Sprintf("%d", colour.red) +
			"," + fmt.Sprintf("%d", colour.green) +
			"," + fmt.Sprintf("%d", colour.blue) + ")"

		if angle == 0 {
			driver.writer.WriteString("<text x=\"" + fmt.Sprintf("%d", x) + "\" y=\"" + fmt.Sprintf("%d", driver.height-y) + "\" " +
				"style=\"" + style + "\" font-family=\"" + driver.fontFamily + "\" font-size=\"" + fmt.Sprintf("%d", driver.fontSize) +
				"\">" + text + "</text>\n")
		} else {
			driver.writer.WriteString("<text x=\"" + fmt.Sprintf("%d", x) + "\" y=\"" + fmt.Sprintf("%d", driver.height-y) + "\" " +
				"transform=\"rotate(" + fmt.Sprintf("%d", angle) + ", " +
				fmt.Sprintf("%d", x) + ", " + fmt.Sprintf("%d", driver.height-y) + ")\" " +
				"style=\"" + style + "\" font-family=\"" + driver.fontFamily + "\" font-size=\"" + fmt.Sprintf("%d", driver.fontSize) +
				"\">" + text + "</text>\n")
		}
	*/

	return nil
}

//	Close finalize the PNG graphic
func (driver *PNG_Driver) Close() error {
	err := png.Encode(driver.writer, driver.image)
	if err != nil {
		return err
	}
	driver.writer.Flush()

	return nil
}
