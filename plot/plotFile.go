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
				dataFileName := match[0][1]
				x_column, err := strconv.Atoi(match[0][2])
				if err != nil {
					return nil, errors.New("x column expected to be numeric: " + err.Error())
				}

				y_column, err := strconv.Atoi(match[0][3])
				if err != nil {
					return nil, errors.New("y column expected to be numeric: " + err.Error())
				}

				//	open the Go-Plot data file and load it
				dataFile, err := os.Open(dataFileName)
				if err != nil {
					return nil, errors.New("fail attempting to open Go-Plot data file: " + err.Error())
				}
				defer dataFile.Close()
				point, err := LoadDataFile(uint8(x_column), uint8(y_column), bufio.NewReader(dataFile))
				if err != nil {
					return nil, errors.New("fail attempting to load Go-Plot data file: " + err.Error())
				}

				plot.set_points = append(plot.set_points, set_points_2d{
					title: "",
					style: POINTS,
					point: point,
				})

				line = ""
				continue
			}
		}
	}

	return plot, nil
}
