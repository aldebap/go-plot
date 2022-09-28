////////////////////////////////////////////////////////////////////////////////
//	plotFile.go  -  Sep-5-2022  -  aldebap
//
//	Parser for Go-Plot files
////////////////////////////////////////////////////////////////////////////////

package plot

import (
	"bufio"
	"errors"
	"fmt"
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
	Style = map[string]uint8{
		"boxes":       BOXES,
		"dots":        DOTS,
		"lines":       LINES,
		"linespoints": LINES_POINTS,
		"points":      POINTS,
	}
)

const (
	DEFAULT_STYLE = "points"
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

	setTerminalRegEx, err := regexp.Compile(`^\s*set\s+terminal\s+(\S+)\s*$`)
	if err != nil {
		return nil, err
	}

	setOutputRegEx, err := regexp.Compile(`^\s*set\s+output\s+"(.+)"\s*$`)
	if err != nil {
		return nil, err
	}

	plotCommandRegEx, err := regexp.Compile(`^\s*plot\s*`)
	if err != nil {
		return nil, err
	}

	dataFileRegEx, err := regexp.Compile(`^\s*"([^"]+)"\s*`)
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

	plotTitleRegEx, err := regexp.Compile(`^\s*title\s+"([^"]+)"\s*`)
	if err != nil {
		return nil, err
	}

	commaSeparatorRegEx, err := regexp.Compile(`^\s*,\s*`)
	if err != nil {
		return nil, err
	}

	//	read the input line by line
	var (
		plot = &Plot_2D{
			Set_points: make([]Set_points_2d, 0),
		}

		line      string
		plotScope bool
	)

	var (
		dataFileName string
		x_column     string = "1"
		y_column     string = "2"
		style        string = DEFAULT_STYLE
		title        string
	)

	for {
		bufLine, isPrefix, err := reader.ReadLine()
		if err != nil {
			break
		}

		line += string(bufLine)

		if !isPrefix {
			//	parse the line using all regex for commands outside plot command scope
			var commandFound bool

			match := setXLabelRegEx.FindAllStringSubmatch(line, -1)
			if len(match) == 1 {
				plot.X_label = match[0][1]
				commandFound = true
			}

			match = setYLabelRegEx.FindAllStringSubmatch(line, -1)
			if len(match) == 1 {
				plot.Y_label = match[0][1]
				commandFound = true
			}

			match = setTerminalRegEx.FindAllStringSubmatch(line, -1)
			if len(match) == 1 {
				var found bool

				plot.Terminal, found = terminal[match[0][1]]
				if !found {
					return nil, errors.New("invalid terminal type: " + match[0][1])
				}
				commandFound = true
			}

			match = setOutputRegEx.FindAllStringSubmatch(line, -1)
			if len(match) == 1 {
				plot.output = match[0][1]
				commandFound = true
			}

			//	if a command was found clean up current line
			if commandFound {
				//	if previous parsed a plot command whose data file was not loaded yet, it's the time for it
				if plotScope && len(dataFileName) > 0 {
					auxSetPoints, err := newSetPoints2D(dataFileName, x_column, y_column, style, title)
					if err != nil {
						return nil, err
					}

					//	erase file name as it was used already
					dataFileName = ""
					x_column = "1"
					y_column = "2"
					style = DEFAULT_STYLE
					title = ""

					plot.Set_points = append(plot.Set_points, *auxSetPoints)
					plot.Set_points[len(plot.Set_points)-1].order = uint8(len(plot.Set_points))

					plotScope = false
				}

				line = ""
				continue
			}

			//	if it's not a configuration command, try to parse plot command options
			for {
				if len(line) == 0 {
					break
				}
				//	fmt.Printf("[debug] line: %s\n", line)

				match = plotCommandRegEx.FindAllStringSubmatch(line, -1)
				if len(match) == 1 {
					plotScope = true

					line = line[len(match[0][0]):]
					continue
				}

				match = dataFileRegEx.FindAllStringSubmatch(line, -1)
				if len(match) == 1 {
					if !plotScope {
						return nil, errors.New("data file specification without a plot command: " + match[0][0])
					}
					dataFileName = match[0][1]

					line = line[len(match[0][0]):]
					continue
				}

				match = dataFilePlotUsingRegEx.FindAllStringSubmatch(line, -1)
				if len(match) == 1 {
					if !plotScope {
						return nil, errors.New("'using' option without a plot command: " + match[0][0])
					}
					x_column = match[0][1]
					y_column = match[0][2]

					line = line[len(match[0][0]):]
					continue
				}

				match = plotWithRegEx.FindAllStringSubmatch(line, -1)
				if len(match) == 1 {
					if !plotScope {
						return nil, errors.New("'with' option without a plot command: " + match[0][0])
					}
					style = match[0][1]

					line = line[len(match[0][0]):]
					continue
				}

				match = plotTitleRegEx.FindAllStringSubmatch(line, -1)
				if len(match) == 1 {
					if !plotScope {
						return nil, errors.New("'title' option without a plot command: " + match[0][0])
					}
					title = match[0][1]

					line = line[len(match[0][0]):]
					continue
				}

				//	when a comma is found in the scope of a plot command, add the data file points
				match = commaSeparatorRegEx.FindAllStringSubmatch(line, -1)
				if len(match) == 1 {
					if !plotScope {
						return nil, errors.New("unexpected syntax: " + match[0][1])
					}
					//	if a data file name was found, add the set of points
					if len(dataFileName) > 0 {
						auxSetPoints, err := newSetPoints2D(dataFileName, x_column, y_column, style, title)
						if err != nil {
							return nil, err
						}

						//	erase file name as it was used already
						dataFileName = ""
						x_column = "1"
						y_column = "2"
						style = DEFAULT_STYLE
						title = ""

						plot.Set_points = append(plot.Set_points, *auxSetPoints)
						plot.Set_points[len(plot.Set_points)-1].order = uint8(len(plot.Set_points))
					}

					line = line[len(match[0][0]):]
					continue
				}
			}
		}
	}

	//	when plot file parsing finishes, if a plot command whose data file was not loaded yet, it's the time for it
	if plotScope && len(dataFileName) > 0 {
		auxSetPoints, err := newSetPoints2D(dataFileName, x_column, y_column, style, title)
		if err != nil {
			return nil, err
		}

		plot.Set_points = append(plot.Set_points, *auxSetPoints)
		plot.Set_points[len(plot.Set_points)-1].order = uint8(len(plot.Set_points))
	}

	return plot, nil
}

//	newSetPoints2D parse string parameters and attempt to create a new set of 2D points
func newSetPoints2D(dataFileName, x_column, y_column, styleDesc, title string) (*Set_points_2d, error) {

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

	num_style, found = Style[styleDesc]
	if !found {
		return nil, errors.New("invalid style: " + styleDesc)
	}

	//	set a default title when necessary
	if len(title) == 0 {
		title = fmt.Sprintf("%s u %d:%d", dataFileName, num_x_column, num_y_column)
	}

	return &Set_points_2d{
		Title: title,
		Style: num_style,
		Point: point,
	}, nil
}
