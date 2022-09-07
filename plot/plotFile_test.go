////////////////////////////////////////////////////////////////////////////////
//	plotFile_test.go  -  Sep-5-2022  -  aldebap
//
//	Test cases for parser for Go-Plot files
////////////////////////////////////////////////////////////////////////////////

package plot

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

//	TestLoadPlotFile unit tests for LoadPlotFile()
func TestLoadPlotFile(t *testing.T) {

	t.Run(">>> LoadPlotFile: xlabel", func(t *testing.T) {
		want := "x axis"

		mockPlotFile := strings.NewReader(`set xlabel "` + want + `"`)
		plot, err := LoadPlotFile(bufio.NewReader(mockPlotFile))
		if err != nil {
			t.Errorf("fail loading plot file: %s", err.Error())
			return
		}

		got := plot.(*Plot_2D).x_label
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: '%s' result: '%s'", want, got)
		}
	})

	t.Run(">>> LoadPlotFile: ylabel", func(t *testing.T) {
		want := "y axis"

		mockPlotFile := strings.NewReader(`set ylabel "` + want + `"`)
		plot, err := LoadPlotFile(bufio.NewReader(mockPlotFile))
		if err != nil {
			t.Errorf("fail loading plot file: %s", err.Error())
			return
		}

		got := plot.(*Plot_2D).y_label
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: '%s' result: '%s'", want, got)
		}
	})

	t.Run(">>> LoadPlotFile: plot \"file\" using x:y", func(t *testing.T) {

		//	create a temporary data file
		tmpDataFile, err := os.CreateTemp("", "goPlotData")
		if err != nil {
			t.Errorf("fail creating plot data file: %s", err.Error())
			return
		}
		defer os.Remove(tmpDataFile.Name())

		_, err = tmpDataFile.Write([]byte("col1 col2 col3\n10 20 30\n40 50 60\n"))
		if err != nil {
			tmpDataFile.Close()
			t.Errorf("fail writing to plot data file: %s", err.Error())
			return
		}
		err = tmpDataFile.Close()
		if err != nil {
			t.Errorf("fail closing the plot data file: %s", err.Error())
			return
		}

		expectedSetPoints := 1
		expectedPoints := 2

		mockPlotFile := strings.NewReader(`plot "` + tmpDataFile.Name() + `" using 1:3`)
		plot, err := LoadPlotFile(bufio.NewReader(mockPlotFile))
		if err != nil {
			t.Errorf("fail loading plot file: %s", err.Error())
			return
		}

		want := expectedSetPoints
		got := len(plot.(*Plot_2D).set_points)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d sets result: %d", want, got)
			return
		}

		want = expectedPoints
		got = len(plot.(*Plot_2D).set_points[0].point)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d points result: %d", want, got)
		}
	})
}
