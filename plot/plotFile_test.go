////////////////////////////////////////////////////////////////////////////////
//	plotFile_test.go  -  Sep-5-2022  -  aldebap
//
//	Test cases for parser for Go-Plot files
////////////////////////////////////////////////////////////////////////////////

package plot

import (
	"bufio"
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
}
