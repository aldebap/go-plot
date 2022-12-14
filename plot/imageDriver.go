////////////////////////////////////////////////////////////////////////////////
//	imageDriver.go  -  Nov-24-2022  -  aldebap
//
//	Implementation of a graphic driver to generate PNG, GIF and JPEG files
////////////////////////////////////////////////////////////////////////////////

package plot

import (
	"bufio"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

type Image_Driver struct {
	writer     *bufio.Writer
	fileFormat string
	width      int64
	height     int64
	image      *image.RGBA
	path       []DriverPoint
	pathColour RGB_colour
	fontFamily string
	fontSize   uint8
	dpi        uint16
	font       *truetype.Font
	ctx        *freetype.Context
}

//	create a new PNG_Driver
func NewPNG_Driver(writer *bufio.Writer) GraphicsDriver {
	const (
		WIDTH       = 640
		HEIGHT      = 480
		FONT_FAMILY = "Verdana"
		FONT_SIZE   = 10
		DPI         = 72
	)

	return &Image_Driver{
		writer:     writer,
		fileFormat: "png",
		width:      WIDTH,
		height:     HEIGHT,
		fontFamily: FONT_FAMILY,
		fontSize:   FONT_SIZE,
		dpi:        DPI,
	}
}

//	create a new GIF
func NewGIF_Driver(writer *bufio.Writer) GraphicsDriver {
	const (
		WIDTH       = 640
		HEIGHT      = 480
		FONT_FAMILY = "Verdana"
		FONT_SIZE   = 10
		DPI         = 72
	)

	return &Image_Driver{
		writer:     writer,
		fileFormat: "gif",
		width:      WIDTH,
		height:     HEIGHT,
		fontFamily: FONT_FAMILY,
		fontSize:   FONT_SIZE,
		dpi:        DPI,
	}
}

//	create a new JPEG
func NewJPEG_Driver(writer *bufio.Writer) GraphicsDriver {
	const (
		WIDTH       = 640
		HEIGHT      = 480
		FONT_FAMILY = "Verdana"
		FONT_SIZE   = 10
		DPI         = 72
	)

	return &Image_Driver{
		writer:     writer,
		fileFormat: "jpeg",
		width:      WIDTH,
		height:     HEIGHT,
		fontFamily: FONT_FAMILY,
		fontSize:   FONT_SIZE,
		dpi:        DPI,
	}
}

//	GetDimensions get the dimensions of the Image graphic
func (driver *Image_Driver) GetDimensions() (width, heigth int64) {
	return driver.width, driver.height
}

//	SetDimensions set the dimensions of the Image graphic
func (driver *Image_Driver) SetDimensions(width int64, height int64) error {
	driver.width = width
	driver.height = height

	driver.image = image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{int(width), int(height)}})

	//	set background as white
	draw.Draw(driver.image, driver.image.Bounds(), image.White, image.ZP, draw.Src)

	return nil
}

//	GetFont get information about the font
func (driver *Image_Driver) GetFont() (fontFamily string, fontSize uint8) {
	return driver.fontFamily, driver.fontSize
}

//	SetFont set information about the font
func (driver *Image_Driver) SetFont(fontFamily string, fontSize uint8) error {
	if driver.image == nil {
		return errors.New("cannot setFont to a non initialized graphics driver")
	}

	driver.fontFamily = fontFamily
	driver.fontSize = fontSize

	//	load font file
	fontBytes, err := ioutil.ReadFile("../res/font/" + driver.fontFamily + ".ttf")
	if err != nil {
		return err
	}

	driver.font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}

	//	create the context for free type fonts
	driver.ctx = freetype.NewContext()

	driver.ctx.SetDPI(float64(driver.dpi))
	driver.ctx.SetFont(driver.font)
	driver.ctx.SetFontSize(float64(driver.fontSize))
	driver.ctx.SetClip(driver.image.Bounds())
	driver.ctx.SetDst(driver.image)
	driver.ctx.SetSrc(image.Black)

	return nil
}

//	Comment write a comment int the the Image graphic
func (driver *Image_Driver) Comment(text string) {
	//	cannot add comments to Image
}

//	Point draws a point in the Image graphic
func (driver *Image_Driver) Point(x, y int64, colour RGB_colour) error {
	if driver.image == nil {
		return errors.New("cannot draw a point to a non initialized graphics driver")
	}

	pointColour := color.RGBA{colour.red, colour.green, colour.blue, 255}
	driver.image.Set(int(x), int(y), pointColour)

	return nil
}

//	Begin a path to draw a connection between a set of points
func (driver *Image_Driver) BeginPath(colour RGB_colour) error {
	driver.path = make([]DriverPoint, 0)
	driver.pathColour = colour

	return nil
}

//	Add a point to a path
func (driver *Image_Driver) PointToPath(x, y int64) error {
	if driver.path == nil {
		return errors.New("cannot add a point to a non initialized path")
	}
	driver.path = append(driver.path, DriverPoint{X: x, Y: y})

	return nil
}

//	End the path
func (driver *Image_Driver) EndPath() error {
	if driver.path == nil {
		return errors.New("cannot end a path not initialized")
	}

	if len(driver.path) < 2 {
		return errors.New("not enough points to draw a path")
	}

	var x1, y1 int64

	for i, point := range driver.path {
		if i == 0 {
			x1 = point.X
			y1 = driver.height - point.Y
		} else {
			err := driver.Line(x1, y1, point.X, driver.height-point.Y, driver.pathColour)
			if err != nil {
				return err
			}

			x1 = point.X
			y1 = driver.height - point.Y
		}
	}

	driver.path = nil

	return nil
}

//	Line draws a line between two points in the Image graphic
func (driver *Image_Driver) Line(x1, y1, x2, y2 int64, colour RGB_colour) error {
	if driver.image == nil {
		return errors.New("cannot draw a line to a non initialized graphics driver")
	}

	lineColour := color.RGBA{colour.red, colour.green, colour.blue, 255}

	//	for better results, change a function variable if angle is bigger the 45 deg
	if y2-y1 > x2-x1 {
		//	swap the points if y2 < y1
		if y2 < y1 {
			aux := x1
			x1 = x2
			x2 = aux

			aux = y1
			y1 = y2
			y2 = aux
		}

		//	use a line math function to draw the line | x: f(y)
		delta := float64(x2-x1) / float64(y2-y1)
		x := float64(x1)

		for y := y1; y <= y2; y++ {
			driver.image.Set(int(x), int(y), lineColour)
			x += delta
		}
	} else {
		//	swap the points if x2 < x1
		if x2 < x1 {
			aux := x1
			x1 = x2
			x2 = aux

			aux = y1
			y1 = y2
			y2 = aux
		}

		//	use a line math function to draw the line | y: f(x)
		delta := float64(y2-y1) / float64(x2-x1)
		y := float64(y1)

		for x := x1; x <= x2; x++ {
			driver.image.Set(int(x), int(y), lineColour)
			y += delta
		}
	}

	return nil
}

//	GetTextBox evaluate the width and height of the rectangle required to draw the text string using a given font size
func (driver *Image_Driver) GetTextBox(text string) (width, height int64) {
	if driver.font == nil {
		return 0, 0
	}

	ttOptions := truetype.Options{
		Size: float64(driver.fontSize),
		DPI:  float64(driver.dpi),
	}
	face := truetype.NewFace(driver.font, &ttOptions)

	//	calculate text width char by char
	width = 0
	for _, char := range text {
		advance, ok := face.GlyphAdvance(char)
		if ok {
			width += int64(advance.Round())
		}
	}
	height = int64(face.Metrics().CapHeight)

	return width, height

}

//	Text writes a string to the specified point in the Image graphic
func (driver *Image_Driver) Text(x, y, angle int64, text string, colour RGB_colour) error {
	if driver.ctx == nil {
		return errors.New("cannot draw text to a non initialized font options")
	}

	var pt = freetype.Pt(int(x), int(driver.height-y)+int(driver.ctx.PointToFixed(float64(driver.fontSize))>>6))
	var textColour = color.RGBA{colour.red, colour.green, colour.blue, 255}
	var err error

	driver.ctx.SetSrc(image.NewUniform(textColour))

	_, err = driver.ctx.DrawString(text, pt)
	if err != nil {
		return err
	}

	return nil
}

//	Close finalize the Image graphic
func (driver *Image_Driver) Close() error {
	switch driver.fileFormat {
	case "png":
		err := png.Encode(driver.writer, driver.image)
		if err != nil {
			return err
		}

	case "gif":
		err := gif.Encode(driver.writer, driver.image, nil)
		if err != nil {
			return err
		}

	case "jpeg":
		err := jpeg.Encode(driver.writer, driver.image, nil)
		if err != nil {
			return err
		}
	}
	driver.writer.Flush()

	return nil
}
