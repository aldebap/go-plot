////////////////////////////////////////////////////////////////////////////////
//	dataFile.go  -  Sep-5-2022  -  aldebap
//
//	Parser for Go-Plot data files
////////////////////////////////////////////////////////////////////////////////

package plot

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//	LoadDataFile load a data file and return a Plot
func LoadDataFile(x_column uint8, y_column uint8, reader *bufio.Reader) ([]point_2d, error) {
	point := make([]point_2d, 0, 10)

	//	read the input line by line
	var line string
	var firstLine = true

	for {
		bufLine, isPrefix, err := reader.ReadLine()
		if err != nil {
			break
		}

		line += string(bufLine)

		if !isPrefix {
			//	ignore first line since it's expected to be a header
			if firstLine {
				firstLine = false
				line = ""
				continue
			}

			column := strings.Split(line, " ")

			//	check if the line have the expected columns
			if len(column) < int(x_column) || len(column) < int(y_column) {
				return nil, errors.New(`line with less columns than expected: "` + line + `"`)
			}

			//	check if the columns are numeric
			x, err := strconv.ParseFloat(column[x_column-1], 64)
			if err != nil {
				return nil, errors.New(`column ` + fmt.Sprintf("%d", x_column) + ` expected to be numeric: "` + line + `"`)
			}

			y, err := strconv.ParseFloat(column[y_column-1], 64)
			if err != nil {
				return nil, errors.New(`column ` + fmt.Sprintf("%d", y_column) + ` expected to be numeric: "` + line + `"`)
			}

			//	add the new point
			point = append(point, point_2d{x: x, y: y})

			line = ""
		}
	}

	return point, nil
}
