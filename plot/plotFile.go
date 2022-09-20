////////////////////////////////////////////////////////////////////////////////
//	plotFile.go  -  Sep-5-2022  -  aldebap
//
//	Parser for Go-Plot files
////////////////////////////////////////////////////////////////////////////////

package plot

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strconv"
)

//	style descriptions for a plot of points
const (
	DESC_DOTS       = "dots"
	DESC_LINES      = "lines"
	DESC_LINES_DOTS = "x"
	DESC_BOXES      = "boxes"
)

//	LoadPlotFile load a plot file and return a Plot
func LoadPlotFile(reader *bufio.Reader) (Plot, error) {
	plot := &Plot_2D{
		set_points: make([]set_points_2d, 0),
	}

	//	compile all regexs required to parse the plot file
	var err error

	setXLabelRegEx, err := regexp.Compile(`^\s*set\s+xlabel\s+"(.+)"\s*$`)
	if err != nil {
		return nil, err
	}

	setYLabelRegEx, err := regexp.Compile(`^\s*set\s+ylabel\s+"(.+)"\s*$`)
	if err != nil {
		return nil, err
	}

	dataFilePlotRegEx, err := regexp.Compile(`^\s*plot\s+"(.+)"\s+using\s+(\d+):(\d+)\s*$`)
	if err != nil {
		return nil, err
	}

	dataFilePlotWithRegEx, err := regexp.Compile(`^\s*plot\s+"(.+)"\s+using\s+(\d+):(\d+)\s*with\s+(\S+)\s*$`)
	if err != nil {
		return nil, err
	}

	//	read the input line by line
	var line string

	for {
		bufLine, isPrefix, err := reader.ReadLine()
		if err != nil {
			break
		}

		line += string(bufLine)

		if !isPrefix {
			//	parse the line using all regex
			match := setXLabelRegEx.FindAllStringSubmatch(line, -1)
			if len(match) == 1 {
				plot.x_label = match[0][1]
				line = ""
				continue
			}

			match = setYLabelRegEx.FindAllStringSubmatch(line, -1)
			if len(match) == 1 {
				plot.y_label = match[0][1]
				line = ""
				continue
			}

			match = dataFilePlotRegEx.FindAllStringSubmatch(line, -1)
			if len(match) == 1 {
				auxSetPoints, err := newSetPoints2D(match[0][1], match[0][2], match[0][3], DESC_DOTS)
				if err != nil {
					return nil, err
				}

				plot.set_points = append(plot.set_points, *auxSetPoints)

				line = ""
				continue
			}

			match = dataFilePlotWithRegEx.FindAllStringSubmatch(line, -1)
			if len(match) == 1 {
				auxSetPoints, err := newSetPoints2D(match[0][1], match[0][2], match[0][3], match[0][4])
				if err != nil {
					return nil, err
				}

				plot.set_points = append(plot.set_points, *auxSetPoints)

				line = ""
				continue
			}
		}
	}

	return plot, nil
}

//	newSetPoints2D parse string parameters and attempt to create a new set of 2D points
func newSetPoints2D(dataFileName, x_column, y_column, style string) (*set_points_2d, error) {

	//	attempt to convert x_column to an int
	num_x_column, err := strconv.Atoi(x_column)
	if err != nil {
		return nil, errors.New("x column expected to be numeric: " + err.Error())
	}

	//	attempt to convert y_column to an int
	num_y_column, err := strconv.Atoi(y_column)
	if err != nil {
		return nil, errors.New("y column expected to be numeric: " + err.Error())
	}

	//	open the Go-Plot data file and load it
	dataFile, err := os.Open(dataFileName)
	if err != nil {
		return nil, errors.New("fail attempting to open Go-Plot data file: " + err.Error())
	}
	defer dataFile.Close()
	point, err := LoadDataFile(uint8(num_x_column), uint8(num_y_column), bufio.NewReader(dataFile))
	if err != nil {
		return nil, errors.New("fail attempting to load Go-Plot data file: " + err.Error())
	}

	//	attempt to convert the style string to an int constant
	var num_style uint8

	switch style {
	case DESC_DOTS:
		num_style = DOTS

	case DESC_LINES:
		num_style = LINES

	case DESC_LINES_DOTS:
		num_style = LINES_DOTS

	case DESC_BOXES:
		num_style = BOXES
	}

	return &set_points_2d{
		title: "",
		style: num_style,
		point: point,
	}, nil
}
