////////////////////////////////////////////////////////////////////////////////
//	plot_2d.go  -  Sep-5-2022  -  aldebap
//
//	Generate a 2D Go-Plot
////////////////////////////////////////////////////////////////////////////////

package plot

import (
	"errors"
	"fmt"
	"math"
)

//	dimensions in pixels for the plot (half of a letter sheet)
const (
	WIDTH  float64 = 640
	HEIGHT float64 = 480
)

//	margins in pixels for the plot
const (
	X_MARGINS float64 = 30
	Y_MARGINS float64 = 30
)

//	styles for a plot of points
const (
	DOTS       uint8 = 1
	LINES      uint8 = 2
	LINES_DOTS uint8 = 3
	BOXES      uint8 = 4
)

const (
	Y_SCALE_DIVISIONS = 7
	SCALE_WIDTH       = 6
	FONT_SIZE         = 10
	POINT_WIDTH       = 8
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
func (p *Plot_2D) GeneratePlot(driver GraphicsDriver) error {

	//	check if there's a plot to be generated
	if len(p.set_points) == 0 {
		return errors.New("no set of points to be plotted")
	}

	if len(p.set_points[0].point) == 0 {
		return errors.New("no points in the first set to be plotted")
	}

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

	//	round the scale to multiples of 10
	min_x = math.Floor(min_x) - float64(int64(min_x)%10)
	min_y = math.Floor(min_y) - float64(int64(min_y)%10)
	max_x = math.Floor(max_x) + float64(10-int64(max_x)%10)
	max_y = math.Floor(max_y) + float64(10-int64(max_y)%10)

	//	set the graphics dimension
	plotWidth := int64(math.Round(WIDTH))
	plotHeight := int64(math.Round(HEIGHT))

	driver.SetDimensions(plotWidth, plotHeight)

	//	generate the plot grid
	GeneratePlotGrid(driver, plotWidth, plotHeight, min_x, min_y, max_x, max_y)

	//	add the X & Y titles
	//	TODO: centralize the text of both titles
	if len(p.x_label) > 0 {
		driver.Text(int64(X_MARGINS)+plotWidth/2, int64(math.Round(Y_MARGINS/4)), 0, FONT_SIZE, p.x_label)
	}
	if len(p.y_label) > 0 {
		driver.Text(int64(math.Round(X_MARGINS/4)), int64(Y_MARGINS)+plotHeight/2, -90, FONT_SIZE, p.y_label)
	}

	//	generate the plot for every set of points
	for _, pointsSet := range p.set_points {
		//	TODO: use different coulours for each set
		pointsSet.GeneratePlot(driver, min_x, min_y, max_x, max_y)
	}

	return nil
}

//	GeneratePlotGrid implementation of 2D Go_Plot grid generation
func GeneratePlotGrid(driver GraphicsDriver, plotWidth, plotHeight int64, min_x, min_y, max_x, max_y float64) {

	//	add the plot grid
	driver.Line(int64(X_MARGINS), int64(Y_MARGINS),
		plotWidth-int64(X_MARGINS), int64(Y_MARGINS))
	driver.Line(int64(X_MARGINS), plotHeight-int64(Y_MARGINS),
		plotWidth-int64(X_MARGINS), plotHeight-int64(Y_MARGINS))

	driver.Line(int64(X_MARGINS), int64(Y_MARGINS),
		int64(X_MARGINS), plotHeight-int64(Y_MARGINS))
	driver.Line(plotWidth-int64(X_MARGINS), int64(Y_MARGINS),
		plotWidth-int64(X_MARGINS), plotHeight-int64(Y_MARGINS))

	//	add the X scale in the plot grid
	xScaleDivisions := int64(math.Ceil(float64(plotWidth/plotHeight) * Y_SCALE_DIVISIONS))

	for i := int64(1); i < int64(xScaleDivisions); i++ {
		x := float64(i * int64(max_x-min_x) / xScaleDivisions)
		scaled_x := int64((float64(plotWidth) - 2*X_MARGINS) * (x - min_x) / (max_x - min_x))

		driver.Line(int64(X_MARGINS)+scaled_x, int64(Y_MARGINS),
			int64(X_MARGINS)+scaled_x, int64(Y_MARGINS)+SCALE_WIDTH)
		driver.Line(int64(X_MARGINS)+scaled_x, plotHeight-int64(Y_MARGINS),
			int64(X_MARGINS)+scaled_x, plotHeight-int64(Y_MARGINS)-SCALE_WIDTH)

		//	TODO: centralize text based on the scale indicator
		driver.Text(int64(X_MARGINS)+scaled_x, int64(Y_MARGINS/2), 0, FONT_SIZE, fmt.Sprintf("%d", int64(x)))
	}

	//	add the Y scale in the plot grid
	yScaleDivisions := int64(Y_SCALE_DIVISIONS)

	for i := int64(1); i < int64(yScaleDivisions); i++ {
		y := float64(i * int64(max_y-min_y) / yScaleDivisions)
		scaled_y := int64((float64(plotHeight) - 2*Y_MARGINS) * (y - min_y) / (max_y - min_y))

		driver.Line(int64(X_MARGINS), int64(Y_MARGINS)+scaled_y,
			int64(X_MARGINS)+SCALE_WIDTH, int64(Y_MARGINS)+scaled_y)
		driver.Line(plotWidth-int64(X_MARGINS), int64(Y_MARGINS)+scaled_y,
			plotWidth-int64(X_MARGINS)-SCALE_WIDTH, int64(Y_MARGINS)+scaled_y)

		//	TODO: centralize text based on the scale indicator
		driver.Text(int64(X_MARGINS/2), int64(Y_MARGINS)+scaled_y, 0, FONT_SIZE, fmt.Sprintf("%d", int64(y)))
	}
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

	//	TODO: implement the other styles
	switch set.style {
	case DOTS:
		//	generate a cross for each point
		for _, point := range set.point {
			scaled_x := (WIDTH - 2*X_MARGINS) * (point.x - min_x) / (max_x - min_x)
			scaled_y := (HEIGHT - 2*Y_MARGINS) * (point.y - min_y) / (max_y - min_y)

			driver.Line(int64(X_MARGINS+scaled_x-POINT_WIDTH/2), int64(Y_MARGINS+scaled_y),
				int64(X_MARGINS+scaled_x+POINT_WIDTH/2), int64(Y_MARGINS+scaled_y))
			driver.Line(int64(X_MARGINS+scaled_x), int64(Y_MARGINS+scaled_y-POINT_WIDTH/2),
				int64(X_MARGINS+scaled_x), int64(Y_MARGINS+scaled_y+POINT_WIDTH/2))
		}

	case LINES:
	case LINES_DOTS:

	case BOXES:
		//	generate a box for each point
		for _, point := range set.point {
			scaled_x1 := (WIDTH - 2*X_MARGINS) * (point.x - 0.5 - min_x) / (max_x - min_x)
			scaled_x2 := (WIDTH - 2*X_MARGINS) * (point.x + 0.5 - min_x) / (max_x - min_x)
			scaled_y1 := (HEIGHT - 2*Y_MARGINS) * (0 - min_y) / (max_y - min_y)
			scaled_y2 := (HEIGHT - 2*Y_MARGINS) * (point.y - min_y) / (max_y - min_y)

			driver.Line(int64(X_MARGINS+scaled_x1), int64(Y_MARGINS+scaled_y1),
				int64(X_MARGINS+scaled_x1), int64(Y_MARGINS+scaled_y2))
			driver.Line(int64(X_MARGINS+scaled_x1), int64(Y_MARGINS+scaled_y2),
				int64(X_MARGINS+scaled_x2), int64(Y_MARGINS+scaled_y2))
			driver.Line(int64(X_MARGINS+scaled_x2), int64(Y_MARGINS+scaled_y1),
				int64(X_MARGINS+scaled_x2), int64(Y_MARGINS+scaled_y2))
		}

	default:
	}

	return nil
}
