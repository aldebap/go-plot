////////////////////////////////////////////////////////////////////////////////
//	plot_2d.go  -  Sep-5-2022  -  aldebap
//
//	Generate a 2D Go-Plot
////////////////////////////////////////////////////////////////////////////////

package plot

import (
	"bufio"
	"errors"
)

//	margins for the plot
const (
	X_MARGINS = 5
	Y_MARGINS = 5
)

//	GeneratePlot implementation of 2D Go_Plot generation
func (p *Plot_2D) GeneratePlot(writer *bufio.Writer) error {

	//	check if there's a plot to be generated
	if len(p.set_points) == 0 {
		return errors.New("no set of points to be plotted")
	}

	if len(p.set_points[0].point) == 0 {
		return errors.New("no points in the first set to be plotted")
	}

	//	evaluate the plot's dimension
	min_x := p.set_points[0].point[0].x
	max_x := min_x
	min_y := p.set_points[0].point[0].y
	max_y := min_y

	for _, pointsSet := range p.set_points {
		for _, point := range pointsSet.point {
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
	}

	width := X_MARGINS + (max_x - min_x) + X_MARGINS
	height := Y_MARGINS + (max_y - min_y) + Y_MARGINS

	//	create the drive and set the graphics dimension
	driver := NewSVG_Driver(writer, int64(width), int64(height))

	_ = driver

	writer.Flush()

	return nil
}
