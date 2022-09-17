////////////////////////////////////////////////////////////////////////////////
//	graphicsDriver.go  -  Sep-8-2022  -  aldebap
//
//	Interface for a general graphics driver
////////////////////////////////////////////////////////////////////////////////

package plot

type GraphicsDriver interface {
	SetDimensions(width int64, height int64) error
	Point(x, y int64) error
	Line(x1, y1, x2, y2 int64) error
	Text(x, y, fontSize int64, text string) error
	Close()
}
