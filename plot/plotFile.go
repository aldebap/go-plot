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

//	terminal descriptions for a plot
var (
	terminal = map[string]uint8{
		"canvas": TERMINAL_CANVAS,
		"png":    TERMINAL_PNG,
		"svg":    TERMINAL_SVG,
	}
)

//	style descriptions for a plot of points
var (
	style = map[string]uint8{
		"dots":  DOTS,
		"lines": LINES,
		"x":     LINES_DOTS,
		"boxes": BOXES,
	}
)

const (
	DEFAULT_STYLE = "dots"
)

//	LoadPlotFile load a plot file and return a Plot
func LoadPlotFile(reader *bufio.Reader) (Plot, error) {

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

	dataFilePlotRegEx, err := regexp.Compile(`^\s*plot\s+"(.+)"\s*`)
	if err != nil {
		return nil, err
	}

	additionalDataFileRegEx, err := regexp.Compile(`^\s*"(.+)"\s*`)
	if err != nil {
		return nil, err
	}

	dataFilePlotUsingRegEx, err := regexp.Compile(`^\s*using\s+(\d+):(\d+)\s*`)
	if err != nil {
		return nil, err
	}

	plotWithRegEx, err := regexp.Compile(`^\s*with\s+([a-z]+)\s*`)
	if err != nil {
		return nil, err
	}

	commaSeparatorRegEx, err := regexp.Compile(`^\s*,\s*$`)
	if err != nil {
		return nil, err
	}

	setTerminalRegEx, err := regexp.Compile(`^\s*set\s+terminal\s+(\S+)\s*$`)
	if err != nil {
		return nil, err
	}

	setOutputRegEx, err := regexp.Compile(`^\s*set\s+output\s+"(.+)"\s*$`)
	if err != nil {
		return nil, err
	}

	//	read the input line by line
	var (
		plot = &Plot_2D{
			set_points: make([]set_points_2d, 0),
		}

		line      string
		plotScope bool
	)

	//	TODO: add a default label for each set of points
	//	TODO: improve the parse the comma as a separator for multiple set of points for a plot
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
				plotScope = false
				continue
			}

			match = setYLabelRegEx.FindAllStringSubmatch(line, -1)
			if len(match) == 1 {
				plot.y_label = match[0][1]
				line = ""
				plotScope = false
				continue
			}

			match = setTerminalRegEx.FindAllStringSubmatch(line, -1)
			if len(match) == 1 {
				var found bool

				plot.terminal, found = terminal[match[0][1]]
				if !found {
					return nil, errors.New("invalid terminal type: " + match[0][1])
				}

				line = ""
				plotScope = false
				continue
			}

			match = setOutputRegEx.FindAllStringSubmatch(line, -1)
			if len(match) == 1 {
				plot.output = match[0][1]
				line = ""
				plotScope = false
				continue
			}

			//	parse all remainig elements
			var (
				dataFileName string
				x_column     string = "1"
				y_column     string = "2"
				style        string = DEFAULT_STYLE
			)

			for {
				if len(line) == 0 {
					//	if a data file name was found, add the set of points
					if len(dataFileName) > 0 {
						auxSetPoints, err := newSetPoints2D(dataFileName, x_column, y_column, style)
						if err != nil {
							return nil, err
						}

						plot.set_points = append(plot.set_points, *auxSetPoints)
					}
					break
				}

				match = dataFilePlotRegEx.FindAllStringSubmatch(line, -1)
				if len(match) == 1 {
					plotScope = true
					dataFileName = match[0][1]

					line = line[len(match[0][0]):]
					continue
				}

				match = additionalDataFileRegEx.FindAllStringSubmatch(line, -1)
				if len(match) == 1 {
					if !plotScope {
						return nil, errors.New("data file specification without a plot command: " + match[0][1])
					}
					dataFileName = match[0][1]

					line = line[len(match[0][0]):]
					continue
				}

				match = dataFilePlotUsingRegEx.FindAllStringSubmatch(line, -1)
				if len(match) == 1 {
					if !plotScope {
						return nil, errors.New("'using' option without a plot command: " + match[0][1])
					}
					x_column = match[0][1]
					y_column = match[0][2]

					line = line[len(match[0][0]):]
					continue
				}

				match = plotWithRegEx.FindAllStringSubmatch(line, -1)
				if len(match) == 1 {
					if !plotScope {
						return nil, errors.New("'with' option without a plot command: " + match[0][1])
					}
					style = match[0][1]

					line = line[len(match[0][0]):]
					continue
				}

				match = commaSeparatorRegEx.FindAllStringSubmatch(line, -1)
				if len(match) == 1 {
					if !plotScope {
						return nil, errors.New("unexpected syntax: " + match[0][1])
					}

					line = line[len(match[0][0]):]
					continue
				}
			}
		}
	}

	return plot, nil
}

//	newSetPoints2D parse string parameters and attempt to create a new set of 2D points
func newSetPoints2D(dataFileName, x_column, y_column, styleDesc string) (*set_points_2d, error) {

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
	var found bool

	num_style, found = style[styleDesc]
	if !found {
		return nil, errors.New("invalid style: " + styleDesc)
	}

	return &set_points_2d{
		title: "",
		style: num_style,
		point: point,
	}, nil
}
