////////////////////////////////////////////////////////////////////////////////
//	graphicsDriver.go  -  Sep-8-2022  -  aldebap
//
//	Interface for a general graphics driver
////////////////////////////////////////////////////////////////////////////////

package plot

type RGB_colour struct {
	red   uint8
	green uint8
	blue  uint8
}

type GraphicsDriver interface {
	SetDimensions(width int64, height int64) error
	Point(x, y int64, colour RGB_colour) error
	Line(x1, y1, x2, y2 int64, colour RGB_colour) error
	Text(x, y, angle, fontSize int64, text string, colour RGB_colour) error
	Close()
}
