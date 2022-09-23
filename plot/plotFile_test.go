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

	t.Run(">>> LoadPlotFile: set xlabel", func(t *testing.T) {
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

	t.Run(">>> LoadPlotFile: set ylabel", func(t *testing.T) {
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

	t.Run(">>> LoadPlotFile: set terminal", func(t *testing.T) {
		want := "canvas"

		mockPlotFile := strings.NewReader(`set terminal ` + want)
		plot, err := LoadPlotFile(bufio.NewReader(mockPlotFile))
		if err != nil {
			t.Errorf("fail loading plot file: %s", err.Error())
			return
		}

		got := plot.(*Plot_2D).terminal
		//	check the result
		if terminal[want] != got {
			t.Errorf("failed parsing plot file: expected: %d (%s) result: %d", terminal[want], want, got)
		}
	})

	t.Run(">>> LoadPlotFile: set terminal (invalid)", func(t *testing.T) {
		want := "invalid terminal type: bmp"

		mockPlotFile := strings.NewReader(`set terminal bmp`)
		_, err := LoadPlotFile(bufio.NewReader(mockPlotFile))
		if err == nil {
			t.Errorf("error expected loading plot file")
			return
		}

		got := err
		//	check the result
		if want != got.Error() {
			t.Errorf("failed parsing plot file: expected error: %s result: %s", want, got)
		}
	})

	t.Run(">>> LoadPlotFile: set output", func(t *testing.T) {
		want := "graphics.jpg"

		mockPlotFile := strings.NewReader(`set output "` + want + `"`)
		plot, err := LoadPlotFile(bufio.NewReader(mockPlotFile))
		if err != nil {
			t.Errorf("fail loading plot file: %s", err.Error())
			return
		}

		got := plot.(*Plot_2D).output
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %s result: %s", want, got)
		}
	})

	t.Run(">>> LoadPlotFile: plot \"file\" (default parameters)", func(t *testing.T) {

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

		mockPlotFile := strings.NewReader(`plot "` + tmpDataFile.Name() + `"`)
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

	t.Run(">>> LoadPlotFile: data file without a plot command", func(t *testing.T) {

		//	create a temporary data file
		tmpDataFile, err := os.CreateTemp("", "goPlotData")
		if err != nil {
			t.Errorf("fail creating plot data file: %s", err.Error())
			return
		}
		defer os.Remove(tmpDataFile.Name())

		mockPlotFile := strings.NewReader(`"` + tmpDataFile.Name() + `" using 1:3`)
		_, err = LoadPlotFile(bufio.NewReader(mockPlotFile))
		if err == nil {
			t.Errorf("error expected loading plot file")
			return
		}

		want := "data file specification without a plot command: " + tmpDataFile.Name()
		got := err
		//	check the result
		if want != got.Error() {
			t.Errorf("failed parsing plot file: expected error: %s result: %s", want, got)
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

	t.Run(">>> LoadPlotFile: plot \"file\" with boxes", func(t *testing.T) {

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
		expectedStyle := "boxes"

		mockPlotFile := strings.NewReader(`plot "` + tmpDataFile.Name() + `" with ` + expectedStyle)
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

		want = int(style[expectedStyle])
		got = int(plot.(*Plot_2D).set_points[0].style)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d (%s) result: %d", want, expectedStyle, got)
		}
	})

	t.Run(">>> LoadPlotFile: plot \"file\" with slopes", func(t *testing.T) {

		//	create a temporary data file
		tmpDataFile, err := os.CreateTemp("", "goPlotData")
		if err != nil {
			t.Errorf("fail creating plot data file: %s", err.Error())
			return
		}
		defer os.Remove(tmpDataFile.Name())

		want := "invalid style: slopes"

		mockPlotFile := strings.NewReader(`plot "` + tmpDataFile.Name() + `" with slopes`)
		_, err = LoadPlotFile(bufio.NewReader(mockPlotFile))
		if err == nil {
			t.Errorf("error expected loading plot file")
			return
		}

		got := err
		//	check the result
		if want != got.Error() {
			t.Errorf("failed parsing plot file: expected error: %s result: %s", want, got)
		}
	})

	t.Run(">>> LoadPlotFile: plot \"file\" using 1:3 with boxes", func(t *testing.T) {

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
		expectedStyle := "boxes"

		mockPlotFile := strings.NewReader(`plot "` + tmpDataFile.Name() + `" using 1:3 with ` + expectedStyle)
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

		want = int(style[expectedStyle])
		got = int(plot.(*Plot_2D).set_points[0].style)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d (%s) result: %d", want, expectedStyle, got)
		}
	})

	t.Run(">>> LoadPlotFile: plot \"file\" using 1:3 with boxes, ...", func(t *testing.T) {

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

		expectedSetPoints := 2
		expectedPoints := 2
		expectedStyle := "boxes"

		mockPlotFile := strings.NewReader(`plot "` + tmpDataFile.Name() + `" using 1:3 with ` + expectedStyle + `, ` +
			`"` + tmpDataFile.Name() + `" using 2:3`)
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

		want = int(style[expectedStyle])
		got = int(plot.(*Plot_2D).set_points[0].style)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d (%s) result: %d", want, expectedStyle, got)
		}
	})
}
