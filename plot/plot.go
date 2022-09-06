////////////////////////////////////////////////////////////////////////////////
//	plot.go  -  Sep-4-2022  -  aldebap
//
//	Data structures for Go-Plot
////////////////////////////////////////////////////////////////////////////////

package plot

import "io"

type Plot interface {
	GeneratePlot(writer *io.Writer) error
}

//	styles for a plot of points
const (
	POINTS       = 1
	LINES        = 2
	LINES_POINTS = 3
	BOXES        = 4
)

//	2D point coordinate
type point_2d struct {
	x float64
	y float64
}

//	2D points list
type set_points_2d struct {
	title string
	style uint8
	point []point_2d
}

//	attributes used to describe a 2D plot
type Plot_2D struct {
	x_label    string
	y_label    string
	set_points []set_points_2d
}
