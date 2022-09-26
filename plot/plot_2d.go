////////////////////////////////////////////////////////////////////////////////
//	plot_2d.go  -  Sep-5-2022  -  aldebap
//
//	Generate a 2D Go-Plot
////////////////////////////////////////////////////////////////////////////////

package plot

import (
	"bufio"
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
	BOXES        uint8 = 1
	DOTS         uint8 = 2
	LINES        uint8 = 3
	LINES_POINTS uint8 = 4
	POINTS       uint8 = 5
)

const (
	Y_SCALE_DIVISIONS  = 7
	SCALE_WIDTH        = 6
	FONT_SIZE          = 10
	POINT_WIDTH        = 8
	COLOUR_TITLE_WIDTH = 10
	TITLE_MARGIN       = 10
)

//	colour pallete for plots
var (
	BLACK = RGB_colour{red: 0, green: 0, blue: 0}
	RED   = RGB_colour{red: 255, green: 0, blue: 0}
	GREEN = RGB_colour{red: 0, green: 255, blue: 0}
	BLUE  = RGB_colour{red: 0, green: 0, blue: 255}

	plotPallete = []RGB_colour{RED,
		GREEN,
		BLUE,
	}
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
	order uint8
}

//	attributes used to describe a 2D plot
type Plot_2D struct {
	x_label    string
	y_label    string
	set_points []set_points_2d
	terminal   uint8
	output     string
}

//	GetOutputFileName return the plot's output file name
func (p *Plot_2D) GetOutputFileName() string {
	return p.output
}

//	GeneratePlot implementation of 2D Go_Plot generation
func (p *Plot_2D) GeneratePlot(plotWriter *bufio.Writer) error {

	//	create the graphics driver
	var driver GraphicsDriver

	switch p.terminal {
	case TERMINAL_CANVAS:
	case TERMINAL_PNG:

	case TERMINAL_SVG:
		driver = NewSVG_Driver(plotWriter)

	default:
		driver = NewSVG_Driver(plotWriter)
	}
	defer driver.Close()

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
	//	TODO: must improve the way to evaluate the scale
	min_x = math.Floor(min_x) - float64(int64(min_x)%10)
	min_y = math.Floor(min_y) - float64(int64(min_y)%10)
	max_x = math.Floor(max_x) + float64(10-int64(max_x)%10)
	max_y = math.Floor(max_y) + float64(10-int64(max_y)%10)

	//	set the graphics dimension
	plotWidth := int64(math.Round(WIDTH))
	plotHeight := int64(math.Round(HEIGHT))

	driver.SetDimensions(plotWidth, plotHeight)

	//	generate the plot grid
	generatePlotGrid(driver, plotWidth, plotHeight, min_x, min_y, max_x, max_y)

	//	add the X & Y titles
	if len(p.x_label) > 0 {
		textWidth, textHeight := driver.GetTextBox(p.x_label, FONT_SIZE)

		driver.Text(int64(X_MARGINS)+plotWidth/2-textWidth/2, int64(Y_MARGINS)-2*SCALE_WIDTH-textHeight, 0, p.x_label, FONT_SIZE, BLACK)
	}
	if len(p.y_label) > 0 {
		textWidth, textHeight := driver.GetTextBox(p.y_label, FONT_SIZE)

		driver.Text(int64(X_MARGINS)-2*SCALE_WIDTH-textHeight, int64(Y_MARGINS)+plotHeight/2-textWidth, -90, p.y_label, FONT_SIZE, BLACK)
	}

	//	generate the plot for every set of points
	for i, pointsSet := range p.set_points {
		pointsSet.generatePlot(driver, plotWidth, plotHeight, min_x, min_y, max_x, max_y, plotPallete[i%len(plotPallete)])
	}

	return nil
}

//	generatePlotGrid implementation of 2D Go_Plot grid generation
func generatePlotGrid(driver GraphicsDriver, plotWidth, plotHeight int64, min_x, min_y, max_x, max_y float64) {

	//	add the plot grid
	driver.Comment("plot grid")
	driver.Line(int64(X_MARGINS), int64(Y_MARGINS),
		plotWidth-int64(X_MARGINS), int64(Y_MARGINS), BLACK)
	driver.Line(int64(X_MARGINS), plotHeight-int64(Y_MARGINS),
		plotWidth-int64(X_MARGINS), plotHeight-int64(Y_MARGINS), BLACK)

	driver.Line(int64(X_MARGINS), int64(Y_MARGINS),
		int64(X_MARGINS), plotHeight-int64(Y_MARGINS), BLACK)
	driver.Line(plotWidth-int64(X_MARGINS), int64(Y_MARGINS),
		plotWidth-int64(X_MARGINS), plotHeight-int64(Y_MARGINS), BLACK)

	//	add the X scale in the plot grid
	driver.Comment("grid x scale")
	xScaleDivisions := int64(math.Ceil(float64(plotWidth/plotHeight) * Y_SCALE_DIVISIONS))

	for i := int64(1); i < int64(xScaleDivisions); i++ {
		x := float64(i * int64(max_x-min_x) / xScaleDivisions)
		scaled_x := int64((float64(plotWidth) - 2*X_MARGINS) * (x - min_x) / (max_x - min_x))

		driver.Line(int64(X_MARGINS)+scaled_x, int64(Y_MARGINS),
			int64(X_MARGINS)+scaled_x, int64(Y_MARGINS)+SCALE_WIDTH, BLACK)
		driver.Line(int64(X_MARGINS)+scaled_x, plotHeight-int64(Y_MARGINS),
			int64(X_MARGINS)+scaled_x, plotHeight-int64(Y_MARGINS)-SCALE_WIDTH, BLACK)

		scaleNumber := fmt.Sprintf("%d", int64(x))
		textWidth, textHeight := driver.GetTextBox(scaleNumber, FONT_SIZE)

		driver.Text(int64(X_MARGINS)+scaled_x-textWidth/2, int64(Y_MARGINS)-SCALE_WIDTH-textHeight, 0, scaleNumber, FONT_SIZE, BLACK)
	}

	//	add the Y scale in the plot grid
	driver.Comment("grid y scale")
	yScaleDivisions := int64(Y_SCALE_DIVISIONS)

	for i := int64(1); i < int64(yScaleDivisions); i++ {
		y := float64(i * int64(max_y-min_y) / yScaleDivisions)
		scaled_y := int64((float64(plotHeight) - 2*Y_MARGINS) * (y - min_y) / (max_y - min_y))

		driver.Line(int64(X_MARGINS), int64(Y_MARGINS)+scaled_y,
			int64(X_MARGINS)+SCALE_WIDTH, int64(Y_MARGINS)+scaled_y, BLACK)
		driver.Line(plotWidth-int64(X_MARGINS), int64(Y_MARGINS)+scaled_y,
			plotWidth-int64(X_MARGINS)-SCALE_WIDTH, int64(Y_MARGINS)+scaled_y, BLACK)

		scaleNumber := fmt.Sprintf("%d", int64(y))
		textWidth, textHeight := driver.GetTextBox(scaleNumber, FONT_SIZE)

		driver.Text(int64(X_MARGINS)-SCALE_WIDTH-textWidth, int64(Y_MARGINS)+scaled_y-textHeight/2, 0, scaleNumber, FONT_SIZE, BLACK)
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
func (set *set_points_2d) generatePlot(driver GraphicsDriver, plotWidth, plotHeight int64, min_x, min_y, max_x, max_y float64, colour RGB_colour) error {
	if len(set.point) == 0 {
		return errors.New("no points in the set")
	}

	driver.Comment("plotting " + set.title)

	switch set.style {
	case BOXES:
		//	TODO: improve the way to draw the boxes
		//	generate a box for each point
		for _, point := range set.point {
			scaled_x1 := (WIDTH - 2*X_MARGINS) * (point.x - 0.5 - min_x) / (max_x - min_x)
			scaled_x2 := (WIDTH - 2*X_MARGINS) * (point.x + 0.5 - min_x) / (max_x - min_x)
			scaled_y1 := (HEIGHT - 2*Y_MARGINS) * (0 - min_y) / (max_y - min_y)
			scaled_y2 := (HEIGHT - 2*Y_MARGINS) * (point.y - min_y) / (max_y - min_y)

			driver.Line(int64(X_MARGINS+scaled_x1), int64(Y_MARGINS+scaled_y1),
				int64(X_MARGINS+scaled_x1), int64(Y_MARGINS+scaled_y2), colour)
			driver.Line(int64(X_MARGINS+scaled_x1), int64(Y_MARGINS+scaled_y2),
				int64(X_MARGINS+scaled_x2), int64(Y_MARGINS+scaled_y2), colour)
			driver.Line(int64(X_MARGINS+scaled_x2), int64(Y_MARGINS+scaled_y1),
				int64(X_MARGINS+scaled_x2), int64(Y_MARGINS+scaled_y2), colour)
		}

	case DOTS:
		//	generate a single dot for each point
		for _, point := range set.point {
			scaled_x := (WIDTH - 2*X_MARGINS) * (point.x - min_x) / (max_x - min_x)
			scaled_y := (HEIGHT - 2*Y_MARGINS) * (point.y - min_y) / (max_y - min_y)

			driver.Point(int64(X_MARGINS+scaled_x), int64(Y_MARGINS+scaled_y), colour)
		}

	case LINES:
		//	generate a line connecting each point
		var prev_scaled_x, prev_scaled_y float64

		for i, point := range set.point {
			scaled_x := (WIDTH - 2*X_MARGINS) * (point.x - min_x) / (max_x - min_x)
			scaled_y := (HEIGHT - 2*Y_MARGINS) * (point.y - min_y) / (max_y - min_y)

			//	in the first iteration, just save the current point
			if i == 0 {
				prev_scaled_x = scaled_x
				prev_scaled_y = scaled_y
				continue
			}

			driver.Line(int64(X_MARGINS+prev_scaled_x), int64(Y_MARGINS+prev_scaled_y),
				int64(X_MARGINS+scaled_x), int64(Y_MARGINS+scaled_y), colour)

			prev_scaled_x = scaled_x
			prev_scaled_y = scaled_y
		}

	case LINES_POINTS:
		//	generate a line connecting each point
		var prev_scaled_x, prev_scaled_y float64

		for i, point := range set.point {
			scaled_x := (WIDTH - 2*X_MARGINS) * (point.x - min_x) / (max_x - min_x)
			scaled_y := (HEIGHT - 2*Y_MARGINS) * (point.y - min_y) / (max_y - min_y)

			//	generate a cross for each point
			driver.Line(int64(X_MARGINS+scaled_x-POINT_WIDTH/2), int64(Y_MARGINS+scaled_y),
				int64(X_MARGINS+scaled_x+POINT_WIDTH/2), int64(Y_MARGINS+scaled_y), colour)
			driver.Line(int64(X_MARGINS+scaled_x), int64(Y_MARGINS+scaled_y-POINT_WIDTH/2),
				int64(X_MARGINS+scaled_x), int64(Y_MARGINS+scaled_y+POINT_WIDTH/2), colour)

			//	in the first iteration, just save the current point
			if i == 0 {
				prev_scaled_x = scaled_x
				prev_scaled_y = scaled_y
				continue
			}

			driver.Line(int64(X_MARGINS+prev_scaled_x), int64(Y_MARGINS+prev_scaled_y),
				int64(X_MARGINS+scaled_x), int64(Y_MARGINS+scaled_y), colour)

			prev_scaled_x = scaled_x
			prev_scaled_y = scaled_y
		}

	case POINTS:
		//	TODO: improve to use a different figure for distinct sets of points
		//	generate a cross for each point
		for _, point := range set.point {
			scaled_x := (WIDTH - 2*X_MARGINS) * (point.x - min_x) / (max_x - min_x)
			scaled_y := (HEIGHT - 2*Y_MARGINS) * (point.y - min_y) / (max_y - min_y)

			driver.Line(int64(X_MARGINS+scaled_x-POINT_WIDTH/2), int64(Y_MARGINS+scaled_y),
				int64(X_MARGINS+scaled_x+POINT_WIDTH/2), int64(Y_MARGINS+scaled_y), colour)
			driver.Line(int64(X_MARGINS+scaled_x), int64(Y_MARGINS+scaled_y-POINT_WIDTH/2),
				int64(X_MARGINS+scaled_x), int64(Y_MARGINS+scaled_y+POINT_WIDTH/2), colour)
		}
	default:
	}

	//	show the title
	textWidth, textHeight := driver.GetTextBox(set.title, FONT_SIZE)

	driver.Line(plotWidth-int64(X_MARGINS)-TITLE_MARGIN-COLOUR_TITLE_WIDTH, plotHeight-int64(Y_MARGINS)-int64(set.order)*(TITLE_MARGIN+textHeight/2),
		plotWidth-int64(X_MARGINS)-TITLE_MARGIN, plotHeight-int64(Y_MARGINS)-int64(set.order)*(TITLE_MARGIN+textHeight/2), colour)

	driver.Text(plotWidth-int64(X_MARGINS)-2*TITLE_MARGIN-COLOUR_TITLE_WIDTH-textWidth,
		plotHeight-int64(Y_MARGINS)-int64(set.order)*(TITLE_MARGIN+textHeight), 0, set.title, FONT_SIZE, BLACK)

	return nil
}
