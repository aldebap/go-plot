////////////////////////////////////////////////////////////////////////////////
//	plotFile.go  -  Sep-5-2022  -  aldebap
//
//	Parser for Go-Plot files
////////////////////////////////////////////////////////////////////////////////

package plot

import (
	"bufio"
	"regexp"
)

//	LoadPlotFile load a plot file and return a Plot
func LoadPlotFile(reader *bufio.Reader) (Plot, error) {
	plot := &Plot_2D{}

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

	//	read the input line by line
	var line string

	//lineReader := bufio.NewReader(*reader)
	for {
		//bufLine, isPrefix, err := lineReader.ReadLine()
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
		}
	}

	return plot, nil
}
