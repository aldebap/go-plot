////////////////////////////////////////////////////////////////////////////////
//	plot_2d.go  -  Sep-5-2022  -  aldebap
//
//	Generate a 2D Go-Plot
////////////////////////////////////////////////////////////////////////////////

package plot

import (
	"bufio"
	"errors"
	"math"
)

//	dimensions in millimiters for the plot (half of a letter sheet)
const (
	WIDTH  float64 = 215.9
	HEIGHT float64 = 279.4 / 2
)

//	margins in millimiters for the plot
const (
	X_MARGINS float64 = 20
	Y_MARGINS float64 = 20
)

//	styles for a plot of points
const (
	POINTS       = 1
	LINES        = 2
	LINES_POINTS = 3
	BOXES        = 4
)

const (
	POINT_WIDTH = 6
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

//	GeneratePlot implementation of 2D Go_Plot generation
func (p *Plot_2D) GeneratePlot(writer *bufio.Writer) error {

	//	check if there's a plot to be generated
	if len(p.set_points) == 0 {
		return errors.New("no set of points to be plotted")
	}

	if len(p.set_points[0].point) == 0 {
		return errors.New("no points in the first set to be plotted")
	}

	//	create the graphics driver
	//	TODO: the output format needs to be decided (or configured!)
	driver := NewSVG_Driver(writer)
	defer driver.Close()

	//	set the graphics dimension
	plotWidth := int64(math.Round(WIDTH))
	plotHeight := int64(math.Round(HEIGHT))

	driver.SetDimensions(plotWidth, plotHeight)

	//	evaluate the data's dimension
	var min_x, min_y, max_x, max_y float64

	min_x, min_y, max_x, max_y, err := p.set_points[0].getMinMax()
	if err != nil {
		return errors.New("error evaluating the min-max of set to be plotted: " + err.Error())
	}

	for _, pointsSet := range p.set_points {
		var set_min_x, set_min_y, set_max_x, set_max_y float64

		set_min_x, set_min_y, set_max_x, set_max_y, err := pointsSet.getMinMax()
		if err != nil {
			return errors.New("error evaluating the min-max of set to be plotted: " + err.Error())
		}

		if set_min_x < min_x {
			min_x = set_min_x
		}
		if set_min_y < min_y {
			min_y = set_min_y
		}
		if set_max_x > max_x {
			max_x = set_max_x
		}
		if set_max_y > max_y {
			max_y = set_max_y
		}
	}

	//	TODO: add the plot grid
	driver.Line(int64(X_MARGINS), int64(Y_MARGINS),
		plotWidth-int64(X_MARGINS), int64(Y_MARGINS))
	driver.Line(int64(X_MARGINS), plotHeight-int64(Y_MARGINS),
		plotWidth-int64(X_MARGINS), plotHeight-int64(Y_MARGINS))

	driver.Line(int64(X_MARGINS), int64(Y_MARGINS),
		int64(X_MARGINS), plotHeight-int64(Y_MARGINS))
	driver.Line(plotWidth-int64(X_MARGINS), int64(Y_MARGINS),
		plotWidth-int64(X_MARGINS), plotHeight-int64(Y_MARGINS))

	//	TODO: add the X & Y titles

	//	generate the plot for every set of points
	for _, pointsSet := range p.set_points {
		pointsSet.GeneratePlot(driver, min_x, min_y, max_x, max_y)
	}

	return nil
}

//	getMinMax get the min-max X & Y values for the points in the set
func (set *set_points_2d) getMinMax() (min_x, min_y, max_x, max_y float64, err error) {

	if len(set.point) == 0 {
		return 0, 0, 0, 0, errors.New("no points in the set")
	}

	//	evaluate the plot's dimension
	min_x = set.point[0].x
	max_x = min_x
	min_y = set.point[0].y
	max_y = min_y

	for _, point := range set.point {
		if point.x < min_x {
			min_x = point.x
		}
		if point.x > max_x {
			max_x = point.x
		}

		if point.y < min_y {
			min_y = point.y
		}
		if point.y > max_y {
			max_y = point.y
		}
	}

	return min_x, min_y, max_x, max_y, nil
}

//	GeneratePlot generate the graphic for the points in the set
func (set *set_points_2d) GeneratePlot(driver GraphicsDriver, min_x, min_y, max_x, max_y float64) error {

	if len(set.point) == 0 {
		return errors.New("no points in the set")
	}

	switch set.style {
	case POINTS:
		//	generate a cross for each point
		for _, point := range set.point {
			scaled_x := (WIDTH - 2*X_MARGINS) * (point.x - min_x) / (max_x - min_x)
			scaled_y := (HEIGHT - 2*Y_MARGINS) * (point.y - min_y) / (max_y - min_y)

			driver.Line(int64(X_MARGINS+scaled_x-POINT_WIDTH/2), int64(Y_MARGINS+scaled_y),
				int64(X_MARGINS+scaled_x+POINT_WIDTH/2), int64(Y_MARGINS+scaled_y))
			driver.Line(int64(X_MARGINS+scaled_x), int64(Y_MARGINS+scaled_y-POINT_WIDTH/2),
				int64(X_MARGINS+scaled_x), int64(Y_MARGINS+scaled_y+POINT_WIDTH/2))
		}

	default:
	}

	return nil
}
