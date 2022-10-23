////////////////////////////////////////////////////////////////////////////////
//	graphicsDriver.go  -  Sep-8-2022  -  aldebap
//
//	Interface for a general graphics driver
////////////////////////////////////////////////////////////////////////////////

package plot

import "fmt"

//	driver point coordinate
type DriverPoint struct {
	X int64
	Y int64
}

type RGB_colour struct {
	red   uint8
	green uint8
	blue  uint8
}

//	Hexa  return the colour as a hexadecimal value
func (c *RGB_colour) Hexa() string {
	return fmt.Sprintf("%02x%02x%02x", c.red, c.green, c.blue)
}

type GraphicsDriver interface {
	GetDimensions() (width, heigth int64)
	SetDimensions(width int64, height int64) error
	GetFont() (fontFamily string, fontSize uint8)
	SetFont(fontFamily string, fontSize uint8)

	Comment(text string)
	Point(x, y int64, colour RGB_colour) error
	BeginPath(colour RGB_colour) error
	PointToPath(x, y int64) error
	EndPath() error
	Line(x1, y1, x2, y2 int64, colour RGB_colour) error
	GetTextBox(text string) (width, height int64)
	Text(x, y, angle int64, text string, colour RGB_colour) error
	Close()
}
