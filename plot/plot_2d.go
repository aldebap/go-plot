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

	"github.com/aldebap/go-plot/expression"
)

//	margins in pixels for the plot
const (
	X_MARGINS float64 = 30
	Y_MARGINS float64 = 30
)

//	styles for a plot of points
const (
	BOXES         uint8 = 1
	DOTS          uint8 = 2
	LINES         uint8 = 3
	LINES_POINTS  uint8 = 4
	POINTS        uint8 = 5
	FUNCTION_PATH uint8 = 6
)

const (
	MIN_X_SCALE_DIVISIONS = 10
	MAX_X_SCALE_DIVISIONS = 20
	MIN_Y_SCALE_DIVISIONS = 10
	MAX_Y_SCALE_DIVISIONS = 20

	SCALE_WIDTH        = 6
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

	plotPallete = []RGB_colour{
		RED,
		GREEN,
		BLUE,
	}
)

//	2D point coordinate
type Point_2d struct {
	X float64
	Y float64
}

//	2D points list
type Set_points_2d struct {
	Title string
	Style uint8
	Point []Point_2d
	order uint8
}

//	2D function
type Function_2d struct {
	Title    string
	Style    uint8
	Function string
	Min_x    float64
	Max_x    float64
	order    uint8
}

//	attributes used to describe a 2D plot
type Plot_2D struct {
	X_label    string
	Y_label    string
	Set_points []Set_points_2d
	Function   []Function_2d
	Terminal   uint8
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

	switch p.Terminal {
	case TERMINAL_CANVAS:
		driver = NewCanvas_Driver(plotWriter)

	case TERMINAL_GIF:
		driver = NewGIF_Driver(plotWriter)

	case TERMINAL_JPEG:
		driver = NewJPEG_Driver(plotWriter)

	case TERMINAL_PNG:
		driver = NewPNG_Driver(plotWriter)

	case TERMINAL_SVG:
		driver = NewSVG_Driver(plotWriter)

	default:
		driver = NewSVG_Driver(plotWriter)
	}
	defer driver.Close()

	//	check if there's a plot to be generated
	if len(p.Set_points) == 0 && len(p.Function) == 0 {
		return errors.New("no set of points or functions to be plotted")
	}

	/*
		//	TODO: improve this validation
		if len(p.Set_points[0].Point) == 0 {
			return errors.New("no points in the first set to be plotted")
		}
	*/

	//	parse every function and generate a set of point for it's graphic
	var function_points []Set_points_2d

	if len(p.Function) > 0 {
		function_points = make([]Set_points_2d, len(p.Function))
		width, _ := driver.GetDimensions()

		for i, function := range p.Function {

			fmt.Printf("[debug] parsing function #%d: %s\n", function.order, function.Function)

			functionExpr, err := expression.NewExpression(function.Function)
			if err != nil {
				return errors.New("error parsing function to be plotted: " + err.Error())
			}

			//	create the symbol table
			symbolTable := expression.NewFloatSymbolTable()

			expression.AddStandardMathFuncs(symbolTable)

			function_points[i].Point = make([]Point_2d, width-2*int64(X_MARGINS)+1)
			function_points[i].Style = FUNCTION_PATH
			function_points[i].Title = function.Title

			for j := 0; j < len(function_points[i].Point); j++ {
				symbolTable.SetValue("x", function.Min_x+float64(j)*(function.Max_x-function.Min_x)/(float64(width)-2*X_MARGINS))

				function_points[i].Point[j].X, err = symbolTable.GetValue("x")
				function_points[i].Point[j].Y, err = functionExpr.Evaluate(symbolTable)
				if err != nil {
					return errors.New("error evaluating function to be plotted: " + err.Error())
				}
			}
		}
	}

	//	evaluate the data's dimension
	var min_x, min_y, max_x, max_y float64
	var err error

	if len(p.Set_points) > 0 {

		min_x, min_y, max_x, max_y, err = p.Set_points[0].getMinMax()
		if err != nil {
			return errors.New("error evaluating the min-max of set to be plotted: " + err.Error())
		}
	} else {

		min_x, min_y, max_x, max_y, err = function_points[0].getMinMax()
		if err != nil {
			return errors.New("error evaluating the min-max of function to be plotted: " + err.Error())
		}
	}

	for _, pointsSet := range p.Set_points {
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

	//	evaluate the function's dimension
	for _, pointsSet := range function_points {
		var set_min_x, set_min_y, set_max_x, set_max_y float64

		set_min_x, set_min_y, set_max_x, set_max_y, err := pointsSet.getMinMax()
		if err != nil {
			return errors.New("error evaluating the min-max of function to be plotted: " + err.Error())
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

	//	round the scale to multiples of 10 when all plots are based on data sets
	if len(p.Function) == 0 {
		min_x = math.Floor(min_x) - float64(int64(math.Floor(min_x))%10)
		min_y = math.Floor(min_y) - float64(int64(math.Floor(min_y))%10)

		max_x = math.Ceil(max_x)
		if int64(max_x)%10 > 0 {
			max_x += float64(10 - int64(max_x)%10)
		}
		max_y = math.Ceil(max_y)
		if int64(max_y)%10 > 0 {
			max_y += float64(10 - int64(max_y)%10)
		}
	}

	//	set the graphics dimension
	width, height := driver.GetDimensions()

	err = driver.SetDimensions(width, height)
	if err != nil {
		return errors.New("error setting plot dimentions: " + err.Error())
	}

	//	set the graphics font
	fontFamily, fontSize := driver.GetFont()

	err = driver.SetFont(fontFamily, fontSize)
	if err != nil {
		return errors.New("error setting plot font: " + err.Error())
	}

	//	generate the plot grid
	p.generatePlotGrid(driver, min_x, min_y, max_x, max_y)

	//	add the X & Y titles
	if len(p.X_label) > 0 {
		textWidth, textHeight := driver.GetTextBox(p.X_label)

		driver.Text(int64(X_MARGINS)+width/2-textWidth/2, int64(Y_MARGINS)-2*SCALE_WIDTH-textHeight, 0, p.X_label, BLACK)
	}
	if len(p.Y_label) > 0 {
		textWidth, textHeight := driver.GetTextBox(p.Y_label)

		driver.Text(int64(X_MARGINS)-2*SCALE_WIDTH-textHeight, int64(Y_MARGINS)+height/2-textWidth, -90, p.Y_label, BLACK)
	}

	//	generate the plot for every set of points
	for i, pointsSet := range p.Set_points {
		pointsSet.generatePlot(driver, width, height, min_x, min_y, max_x, max_y, plotPallete[i%len(plotPallete)])
	}

	//	generate the plot for every function
	for i, pointsSet := range function_points {
		pointsSet.generatePlot(driver, width, height, min_x, min_y, max_x, max_y, plotPallete[i%len(plotPallete)])
	}

	return nil
}

//	generatePlotGrid implementation of 2D Go_Plot grid generation
func (p *Plot_2D) generatePlotGrid(driver GraphicsDriver, min_x, min_y, max_x, max_y float64) {

	fmt.Printf("[debug] min (%f, %f) max (%f, %f)\n", min_x, min_y, max_x, max_y)

	width, height := driver.GetDimensions()

	//	add the plot grid
	driver.Comment("plot grid")

	driver.Line(int64(X_MARGINS), int64(Y_MARGINS),
		width-int64(X_MARGINS), int64(Y_MARGINS), BLACK)
	driver.Line(int64(X_MARGINS), height-int64(Y_MARGINS),
		width-int64(X_MARGINS), height-int64(Y_MARGINS), BLACK)

	driver.Line(int64(X_MARGINS), int64(Y_MARGINS),
		int64(X_MARGINS), height-int64(Y_MARGINS), BLACK)
	driver.Line(width-int64(X_MARGINS), int64(Y_MARGINS),
		width-int64(X_MARGINS), height-int64(Y_MARGINS), BLACK)

	//	add the X scale in the plot grid
	driver.Comment("grid x scale")

	xScaleDivisions := int64(max_x-min_x) / 10

	if xScaleDivisions < MIN_X_SCALE_DIVISIONS {
		xScaleDivisions = MIN_X_SCALE_DIVISIONS
	} else {
		if xScaleDivisions > MAX_X_SCALE_DIVISIONS {
			xScaleDivisions = MAX_X_SCALE_DIVISIONS
		}
	}

	for i := int64(0); i <= int64(xScaleDivisions); i++ {
		x := min_x + float64(i*int64(max_x-min_x)/xScaleDivisions)
		scaled_x := int64((float64(width) - 2*X_MARGINS) * (x - min_x) / (max_x - min_x))

		driver.Line(int64(X_MARGINS)+scaled_x, int64(Y_MARGINS),
			int64(X_MARGINS)+scaled_x, int64(Y_MARGINS)+SCALE_WIDTH, BLACK)
		driver.Line(int64(X_MARGINS)+scaled_x, height-int64(Y_MARGINS),
			int64(X_MARGINS)+scaled_x, height-int64(Y_MARGINS)-SCALE_WIDTH, BLACK)

		scaleNumber := fmt.Sprintf("%d", int64(x))
		textWidth, textHeight := driver.GetTextBox(scaleNumber)

		driver.Text(int64(X_MARGINS)+scaled_x-textWidth/2, int64(Y_MARGINS)-SCALE_WIDTH-textHeight, 0, scaleNumber, BLACK)
	}

	//	TODO: there's a bug here !
	//	add the Y scale in the plot grid
	driver.Comment("grid y scale")

	yScaleDivisions := int64(max_y-min_y) / 10

	if yScaleDivisions < MIN_Y_SCALE_DIVISIONS {
		yScaleDivisions = MIN_Y_SCALE_DIVISIONS
	} else {
		if yScaleDivisions > MAX_Y_SCALE_DIVISIONS {
			yScaleDivisions = MAX_Y_SCALE_DIVISIONS
		}
	}

	for i := int64(0); i <= int64(yScaleDivisions); i++ {
		y := min_y + float64(i*int64(max_y-min_y)/yScaleDivisions)
		scaled_y := int64((float64(height) - 2*Y_MARGINS) * (y - min_y) / (max_y - min_y))

		driver.Line(int64(X_MARGINS), int64(Y_MARGINS)+scaled_y,
			int64(X_MARGINS)+SCALE_WIDTH, int64(Y_MARGINS)+scaled_y, BLACK)
		driver.Line(width-int64(X_MARGINS), int64(Y_MARGINS)+scaled_y,
			width-int64(X_MARGINS)-SCALE_WIDTH, int64(Y_MARGINS)+scaled_y, BLACK)

		scaleNumber := fmt.Sprintf("%d", int64(y))
		textWidth, textHeight := driver.GetTextBox(scaleNumber)

		driver.Text(int64(X_MARGINS)-SCALE_WIDTH-textWidth, int64(Y_MARGINS)+scaled_y-textHeight/2, 0, scaleNumber, BLACK)
	}
}

//	getMinMax get the min-max X & Y values for the points in the set
func (set *Set_points_2d) getMinMax() (min_x, min_y, max_x, max_y float64, err error) {

	if len(set.Point) == 0 {
		return 0, 0, 0, 0, errors.New("no points in the set")
	}

	//	evaluate the plot's dimension
	min_x = set.Point[0].X
	max_x = min_x
	min_y = set.Point[0].Y
	max_y = min_y

	for _, point := range set.Point {
		if point.X < min_x {
			min_x = point.X
		}
		if point.X > max_x {
			max_x = point.X
		}

		if point.Y < min_y {
			min_y = point.Y
		}
		if point.Y > max_y {
			max_y = point.Y
		}
	}

	//	when the style is "boxes", add left and right margins
	if set.Style == BOXES {
		var meanXInterval float64

		for i, _ := range set.Point {
			if i == 0 {
				continue
			}
			meanXInterval += set.Point[i].X - set.Point[i-1].X
		}
		meanXInterval /= float64(len(set.Point) - 1)

		min_x -= 3 * meanXInterval
		max_x += 3 * meanXInterval
	}

	return min_x, min_y, max_x, max_y, nil
}

//	GeneratePlot generate the graphic for the points in the set
func (set *Set_points_2d) generatePlot(driver GraphicsDriver, plotWidth, plotHeight int64, min_x, min_y, max_x, max_y float64, colour RGB_colour) error {

	if len(set.Point) == 0 {
		return errors.New("no points in the set")
	}

	driver.Comment("plotting " + set.Title)

	switch set.Style {
	case BOXES:
		//	get the mean interval between consecutive pairs of x points
		var meanXInterval float64

		for i, _ := range set.Point {
			if i == 0 {
				continue
			}
			meanXInterval += set.Point[i].X - set.Point[i-1].X
		}
		meanXInterval /= float64(len(set.Point) - 1)
		halfBoxWidth := meanXInterval / 2

		//	generate an open box for each point
		var scaled_x1, scaled_x2, scaled_y1, scaled_y2 float64
		var previousScaled_y2 float64

		scaled_y1 = (float64(plotHeight) - 2*Y_MARGINS) * (0 - min_y) / (max_y - min_y)

		for _, point := range set.Point {

			scaled_x1 = (float64(plotWidth) - 2*X_MARGINS) * (point.X - halfBoxWidth - min_x) / (max_x - min_x)
			scaled_x2 = (float64(plotWidth) - 2*X_MARGINS) * (point.X + halfBoxWidth - min_x) / (max_x - min_x)
			scaled_y2 = (float64(plotHeight) - 2*Y_MARGINS) * (point.Y - min_y) / (max_y - min_y)

			if previousScaled_y2 <= scaled_y2 {
				driver.Line(int64(X_MARGINS+scaled_x1), int64(Y_MARGINS+scaled_y1),
					int64(X_MARGINS+scaled_x1), int64(Y_MARGINS+scaled_y2), colour)
			} else {
				driver.Line(int64(X_MARGINS+scaled_x1), int64(Y_MARGINS+scaled_y1),
					int64(X_MARGINS+scaled_x1), int64(Y_MARGINS+previousScaled_y2), colour)
			}
			driver.Line(int64(X_MARGINS+scaled_x1), int64(Y_MARGINS+scaled_y2),
				int64(X_MARGINS+scaled_x2), int64(Y_MARGINS+scaled_y2), colour)
			driver.Line(int64(X_MARGINS+scaled_x2), int64(Y_MARGINS+scaled_y1),
				int64(X_MARGINS+scaled_x2), int64(Y_MARGINS+scaled_y2), colour)

			previousScaled_y2 = scaled_y2
		}

		//	close the last box
		driver.Line(int64(X_MARGINS+scaled_x2), int64(Y_MARGINS+scaled_y1),
			int64(X_MARGINS+scaled_x2), int64(Y_MARGINS+scaled_y2), colour)

	case DOTS:
		//	generate a single dot for each point
		for _, point := range set.Point {
			scaled_x := (float64(plotWidth) - 2*X_MARGINS) * (point.X - min_x) / (max_x - min_x)
			scaled_y := (float64(plotHeight) - 2*Y_MARGINS) * (point.Y - min_y) / (max_y - min_y)

			driver.Point(int64(X_MARGINS+scaled_x), int64(Y_MARGINS+scaled_y), colour)
		}

	case LINES:
		//	generate a line connecting each point
		var prev_scaled_x, prev_scaled_y float64

		for i, point := range set.Point {
			scaled_x := (float64(plotWidth) - 2*X_MARGINS) * (point.X - min_x) / (max_x - min_x)
			scaled_y := (float64(plotHeight) - 2*Y_MARGINS) * (point.Y - min_y) / (max_y - min_y)

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

		for i, point := range set.Point {
			scaled_x := (float64(plotWidth) - 2*X_MARGINS) * (point.X - min_x) / (max_x - min_x)
			scaled_y := (float64(plotHeight) - 2*Y_MARGINS) * (point.Y - min_y) / (max_y - min_y)

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
		for _, point := range set.Point {
			scaled_x := (float64(plotWidth) - 2*X_MARGINS) * (point.X - min_x) / (max_x - min_x)
			scaled_y := (float64(plotHeight) - 2*Y_MARGINS) * (point.Y - min_y) / (max_y - min_y)

			driver.Line(int64(X_MARGINS+scaled_x-POINT_WIDTH/2), int64(Y_MARGINS+scaled_y),
				int64(X_MARGINS+scaled_x+POINT_WIDTH/2), int64(Y_MARGINS+scaled_y), colour)
			driver.Line(int64(X_MARGINS+scaled_x), int64(Y_MARGINS+scaled_y-POINT_WIDTH/2),
				int64(X_MARGINS+scaled_x), int64(Y_MARGINS+scaled_y+POINT_WIDTH/2), colour)
		}

	case FUNCTION_PATH:
		//	generate a path connecting each point
		driver.BeginPath(colour)

		for _, point := range set.Point {
			scaled_x := (float64(plotWidth) - 2*X_MARGINS) * (point.X - min_x) / (max_x - min_x)
			scaled_y := (float64(plotHeight) - 2*Y_MARGINS) * (point.Y - min_y) / (max_y - min_y)

			driver.PointToPath(int64(X_MARGINS+scaled_x), int64(Y_MARGINS+scaled_y))
		}
		driver.EndPath()

	default:
	}

	//	show the title
	textWidth, textHeight := driver.GetTextBox(set.Title)

	driver.Line(plotWidth-int64(X_MARGINS)-TITLE_MARGIN-COLOUR_TITLE_WIDTH, plotHeight-int64(Y_MARGINS)-int64(set.order)*(TITLE_MARGIN+textHeight/2),
		plotWidth-int64(X_MARGINS)-TITLE_MARGIN, plotHeight-int64(Y_MARGINS)-int64(set.order)*(TITLE_MARGIN+textHeight/2), colour)

	driver.Text(plotWidth-int64(X_MARGINS)-2*TITLE_MARGIN-COLOUR_TITLE_WIDTH-textWidth,
		plotHeight-int64(Y_MARGINS)-int64(set.order)*(TITLE_MARGIN+textHeight), 0, set.Title, BLACK)

	return nil
}
