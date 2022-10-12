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

		got := plot.(*Plot_2D).X_label
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

		got := plot.(*Plot_2D).Y_label
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

		got := plot.(*Plot_2D).Terminal
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

	t.Run(">>> LoadPlotFile: plot function (default parameters)", func(t *testing.T) {

		expectedSetPoints := 0
		expectedFunctions := 1

		mockPlotFile := strings.NewReader(`plot sin(x)`)
		plot, err := LoadPlotFile(bufio.NewReader(mockPlotFile))
		if err != nil {
			t.Errorf("fail loading plot file: %s", err.Error())
			return
		}

		want := expectedSetPoints
		got := len(plot.(*Plot_2D).Set_points)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d sets result: %d", want, got)
			return
		}

		want = expectedFunctions
		got = len(plot.(*Plot_2D).Function)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d functions result: %d", want, got)
		}
	})

	t.Run(">>> LoadPlotFile: function without a plot command", func(t *testing.T) {

		mockPlotFile := strings.NewReader(`sin(x)`)
		_, err := LoadPlotFile(bufio.NewReader(mockPlotFile))
		if err == nil {
			t.Errorf("error expected loading plot file")
			return
		}

		want := `function specification without a plot command: sin(x)`
		got := strings.TrimRight(err.Error(), " ")
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected error: %s result: %s", want, got)
		}
	})

	t.Run(">>> LoadPlotFile: plot function with interval", func(t *testing.T) {

		expectedSetPoints := 0
		expectedFunctions := 1

		mockPlotFile := strings.NewReader(`plot [0:3.14] sin(x)`)
		plot, err := LoadPlotFile(bufio.NewReader(mockPlotFile))
		if err != nil {
			t.Errorf("fail loading plot file: %s", err.Error())
			return
		}

		want := expectedSetPoints
		got := len(plot.(*Plot_2D).Set_points)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d sets result: %d", want, got)
			return
		}

		want = expectedFunctions
		got = len(plot.(*Plot_2D).Function)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d functions result: %d", want, got)
		}
	})

	t.Run(">>> LoadPlotFile: plot function with signed interval numbers", func(t *testing.T) {

		expectedSetPoints := 0
		expectedFunctions := 1

		mockPlotFile := strings.NewReader(`plot [-3.14:+3.14] sin(x)`)
		plot, err := LoadPlotFile(bufio.NewReader(mockPlotFile))
		if err != nil {
			t.Errorf("fail loading plot file: %s", err.Error())
			return
		}

		want := expectedSetPoints
		got := len(plot.(*Plot_2D).Set_points)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d sets result: %d", want, got)
			return
		}

		want = expectedFunctions
		got = len(plot.(*Plot_2D).Function)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d functions result: %d", want, got)
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
		got := len(plot.(*Plot_2D).Set_points)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d sets result: %d", want, got)
			return
		}

		want = expectedPoints
		got = len(plot.(*Plot_2D).Set_points[0].Point)
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

		want := `data file specification without a plot command: "` + tmpDataFile.Name() + `"`
		got := strings.TrimRight(err.Error(), " ")
		//	check the result
		if want != got {
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
		got := len(plot.(*Plot_2D).Set_points)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d sets result: %d", want, got)
			return
		}

		want = expectedPoints
		got = len(plot.(*Plot_2D).Set_points[0].Point)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d points result: %d", want, got)
		}
	})

	t.Run(">>> LoadPlotFile: using option without a plot command", func(t *testing.T) {

		mockPlotFile := strings.NewReader(`using 1:3`)
		_, err := LoadPlotFile(bufio.NewReader(mockPlotFile))
		if err == nil {
			t.Errorf("error expected loading plot file")
			return
		}

		want := "'using' option without a plot command: using 1:3"
		got := err
		//	check the result
		if want != got.Error() {
			t.Errorf("failed parsing plot file: expected error: %s result: %s", want, got)
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
		got := len(plot.(*Plot_2D).Set_points)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d sets result: %d", want, got)
			return
		}

		want = expectedPoints
		got = len(plot.(*Plot_2D).Set_points[0].Point)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d points result: %d", want, got)
		}

		want = int(Style[expectedStyle])
		got = int(plot.(*Plot_2D).Set_points[0].Style)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d (%s) result: %d", want, expectedStyle, got)
		}
	})

	t.Run(">>> LoadPlotFile: with option without a plot command", func(t *testing.T) {

		mockPlotFile := strings.NewReader(`with boxes`)
		_, err := LoadPlotFile(bufio.NewReader(mockPlotFile))
		if err == nil {
			t.Errorf("error expected loading plot file")
			return
		}

		want := "'with' option without a plot command: with boxes"
		got := err
		//	check the result
		if want != got.Error() {
			t.Errorf("failed parsing plot file: expected error: %s result: %s", want, got)
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

	t.Run(">>> LoadPlotFile: plot \"file\" title \"description\"", func(t *testing.T) {

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
		expectedTitle := "annual inflation"

		mockPlotFile := strings.NewReader(`plot "` + tmpDataFile.Name() + `" title "` + expectedTitle + `"`)
		plot, err := LoadPlotFile(bufio.NewReader(mockPlotFile))
		if err != nil {
			t.Errorf("fail loading plot file: %s", err.Error())
			return
		}

		want := expectedSetPoints
		got := len(plot.(*Plot_2D).Set_points)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d sets result: %d", want, got)
			return
		}

		want = expectedPoints
		got = len(plot.(*Plot_2D).Set_points[0].Point)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d points result: %d", want, got)
		}

		wantString := expectedTitle
		gotString := plot.(*Plot_2D).Set_points[0].Title
		//	check the result
		if wantString != gotString {
			t.Errorf("failed parsing plot file: expected: %s result: %s", wantString, gotString)
		}
	})

	t.Run(">>> LoadPlotFile: title option without a plot command", func(t *testing.T) {

		mockPlotFile := strings.NewReader(`title "outside a plot"`)
		_, err := LoadPlotFile(bufio.NewReader(mockPlotFile))
		if err == nil {
			t.Errorf("error expected loading plot file")
			return
		}

		want := `'title' option without a plot command: title "outside a plot"`
		got := err
		//	check the result
		if want != got.Error() {
			t.Errorf("failed parsing plot file: expected error: %s result: %s", want, got)
		}
	})

	t.Run(">>> LoadPlotFile: plot \"file\" using 1:3 with boxes (multiple options)", func(t *testing.T) {

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
		got := len(plot.(*Plot_2D).Set_points)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d sets result: %d", want, got)
			return
		}

		want = expectedPoints
		got = len(plot.(*Plot_2D).Set_points[0].Point)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d points result: %d", want, got)
		}

		want = int(Style[expectedStyle])
		got = int(plot.(*Plot_2D).Set_points[0].Style)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d (%s) result: %d", want, expectedStyle, got)
		}
	})

	t.Run(">>> LoadPlotFile: plot \"file\" using 1:3 with boxes (options in multiple lines)", func(t *testing.T) {

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

		mockPlotFile := strings.NewReader("plot \"" + tmpDataFile.Name() + "\"\nusing 1:3\nwith " + expectedStyle)
		plot, err := LoadPlotFile(bufio.NewReader(mockPlotFile))
		if err != nil {
			t.Errorf("fail loading plot file: %s", err.Error())
			return
		}

		want := expectedSetPoints
		got := len(plot.(*Plot_2D).Set_points)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d sets result: %d", want, got)
			return
		}

		want = expectedPoints
		got = len(plot.(*Plot_2D).Set_points[0].Point)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d points result: %d", want, got)
		}

		want = int(Style[expectedStyle])
		got = int(plot.(*Plot_2D).Set_points[0].Style)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d (%s) result: %d", want, expectedStyle, got)
		}
	})

	t.Run(">>> LoadPlotFile: plot \"file\" using 1:3 with boxes (followed by a global command)", func(t *testing.T) {

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
		expectedTerminal := "canvas"

		mockPlotFile := strings.NewReader("plot \"" + tmpDataFile.Name() + "\"\nusing 1:3\nwith " + expectedStyle + "\n" +
			"set terminal " + expectedTerminal)
		plot, err := LoadPlotFile(bufio.NewReader(mockPlotFile))
		if err != nil {
			t.Errorf("fail loading plot file: %s", err.Error())
			return
		}

		want := expectedSetPoints
		got := len(plot.(*Plot_2D).Set_points)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d sets result: %d", want, got)
			return
		}

		want = expectedPoints
		got = len(plot.(*Plot_2D).Set_points[0].Point)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d points result: %d", want, got)
		}

		want = int(Style[expectedStyle])
		got = int(plot.(*Plot_2D).Set_points[0].Style)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d (%s) result: %d", want, expectedStyle, got)
		}

		wantString := expectedTerminal
		gotString := plot.(*Plot_2D).Terminal
		//	check the result
		if terminal[wantString] != gotString {
			t.Errorf("failed parsing plot file: expected: %d (%s) result: %d", terminal[wantString], wantString, gotString)
		}
	})

	t.Run(">>> LoadPlotFile: plot both function and data file", func(t *testing.T) {

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

		expectedFunctions := 1
		expectedFunctionStyle := "lines"
		expectedSetPoints := 1
		expectedTitle := "annual inflation"

		mockPlotFile := strings.NewReader(`plot [0:+3.14] sin(x) with ` + expectedFunctionStyle + `, ` +
			`"` + tmpDataFile.Name() + `" using 2:3 title "` + expectedTitle + `"`)
		plot, err := LoadPlotFile(bufio.NewReader(mockPlotFile))
		if err != nil {
			t.Errorf("fail loading plot file: %s", err.Error())
			return
		}

		want := expectedFunctions
		got := len(plot.(*Plot_2D).Function)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d functions result: %d", want, got)
		}

		want = int(Style[expectedFunctionStyle])
		got = int(plot.(*Plot_2D).Function[0].Style)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d (%s) result: %d", want, expectedFunctionStyle, got)
		}

		want = expectedSetPoints
		got = len(plot.(*Plot_2D).Set_points)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d sets result: %d", want, got)
			return
		}

		wantString := expectedTitle
		gotString := plot.(*Plot_2D).Set_points[0].Title
		//	check the result
		if wantString != gotString {
			t.Errorf("failed parsing plot file: expected: %s result: %s", wantString, gotString)
		}
	})

	t.Run(">>> LoadPlotFile: plot both function and data file in same description", func(t *testing.T) {

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

		mockPlotFile := strings.NewReader(`plot [0:+3.14] sin(x) "` + tmpDataFile.Name() + `" using 2:3`)
		_, err = LoadPlotFile(bufio.NewReader(mockPlotFile))
		if err == nil {
			t.Errorf("error expected loading plot file")
			return
		}

		want := `function and data file must be described separate in plot command`
		got := err
		//	check the result
		if want != got.Error() {
			t.Errorf("failed parsing plot file: expected error: %s result: %s", want, got)
		}
	})

	t.Run(">>> LoadPlotFile: plot \"file\" using 1:3 with boxes, ...(two data files)", func(t *testing.T) {

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
		expectedTitle := "annual inflation"

		mockPlotFile := strings.NewReader(`plot "` + tmpDataFile.Name() + `" using 1:3 with ` + expectedStyle + `, ` +
			`"` + tmpDataFile.Name() + `" using 2:3 title "` + expectedTitle + `"`)
		plot, err := LoadPlotFile(bufio.NewReader(mockPlotFile))
		if err != nil {
			t.Errorf("fail loading plot file: %s", err.Error())
			return
		}

		want := expectedSetPoints
		got := len(plot.(*Plot_2D).Set_points)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d sets result: %d", want, got)
			return
		}

		want = expectedPoints
		got = len(plot.(*Plot_2D).Set_points[0].Point)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d points result: %d", want, got)
		}

		want = int(Style[expectedStyle])
		got = int(plot.(*Plot_2D).Set_points[0].Style)
		//	check the result
		if want != got {
			t.Errorf("failed parsing plot file: expected: %d (%s) result: %d", want, expectedStyle, got)
		}

		wantString := expectedTitle
		gotString := plot.(*Plot_2D).Set_points[1].Title
		//	check the result
		if wantString != gotString {
			t.Errorf("failed parsing plot file: expected: %s result: %s", wantString, gotString)
		}
	})
}

//	TestNewFunction2D unit tests for newFunction2D()
func TestNewFunction2D(t *testing.T) {

	t.Run(">>> newFunction2D: non numerical min_x", func(t *testing.T) {
		want := `min x expected to be numeric: strconv.ParseFloat: parsing "one": invalid syntax`

		_, got := newFunction2D("sin(x)", "one", "2", "boxes", "")
		//	check the result
		if want != got.Error() {
			t.Errorf("failed creating a new function 2D: expected error: %s result: %s", want, got)
		}
	})

	t.Run(">>> newFunction2D: non numerical max_x", func(t *testing.T) {
		want := `max x expected to be numeric: strconv.ParseFloat: parsing "two": invalid syntax`

		_, got := newFunction2D("sin(x)", "1", "two", "boxes", "")
		//	check the result
		if want != got.Error() {
			t.Errorf("failed creating a new function 2D: expected error: %s result: %s", want, got)
		}
	})

	t.Run(">>> newFunction2D: invalid style", func(t *testing.T) {

		want := `invalid style: circles`

		_, got := newFunction2D("sin(x)", "-10", "+10", "circles", "")
		//	check the result
		if want != got.Error() {
			t.Errorf("failed creating a new function 2D: expected error: %s result: %s", want, got)
		}
	})

	t.Run(">>> newFunction2D: valid scenario (default title)", func(t *testing.T) {

		expectedStyle := "boxes"
		expectedTitle := "sin(x)"

		function, err := newFunction2D("sin(x)", "-10", "+10", expectedStyle, "")
		if err != nil {
			t.Errorf("failed creating a new function 2D: %s", err.Error())
			return
		}

		want := int(Style[expectedStyle])
		got := int(function.Style)
		//	check the result
		if want != got {
			t.Errorf("failed creating a new function: expected: %d (%s) result: %d", want, expectedStyle, got)
		}

		wantString := expectedTitle
		gotString := function.Title
		//	check the result
		if wantString != gotString {
			t.Errorf("failed creating a new function: expected: %s result: %s", wantString, gotString)
		}
	})

	t.Run(">>> newFunction2D: valid scenario", func(t *testing.T) {

		expectedStyle := "boxes"
		expectedTitle := "senoid"

		function, err := newFunction2D("sin(x)", "-10", "+10", expectedStyle, expectedTitle)
		if err != nil {
			t.Errorf("failed creating a new function 2D: %s", err.Error())
			return
		}

		want := int(Style[expectedStyle])
		got := int(function.Style)
		//	check the result
		if want != got {
			t.Errorf("failed creating a new function: expected: %d (%s) result: %d", want, expectedStyle, got)
		}

		wantString := expectedTitle
		gotString := function.Title
		//	check the result
		if wantString != gotString {
			t.Errorf("failed creating a new function: expected: %s result: %s", wantString, gotString)
		}
	})
}

//	TestNewSetPoints2D unit tests for newSetPoints2D()
func TestNewSetPoints2D(t *testing.T) {

	t.Run(">>> newSetPoints2D: non numerical x_column", func(t *testing.T) {
		want := `x column expected to be numeric: strconv.Atoi: parsing "one": invalid syntax`

		_, got := newSetPoints2D("fake", "one", "2", "boxes", "")
		//	check the result
		if want != got.Error() {
			t.Errorf("failed creating a new set of points: expected error: %s result: %s", want, got)
		}
	})

	t.Run(">>> newSetPoints2D: non numerical y_column", func(t *testing.T) {
		want := `y column expected to be numeric: strconv.Atoi: parsing "two": invalid syntax`

		_, got := newSetPoints2D("fake", "1", "two", "boxes", "")
		//	check the result
		if want != got.Error() {
			t.Errorf("failed creating a new set of points: expected error: %s result: %s", want, got)
		}
	})

	t.Run(">>> newSetPoints2D: invalid data file name", func(t *testing.T) {
		want := `fail attempting to open Go-Plot data file: open fake: no such file or directory`

		_, got := newSetPoints2D("fake", "1", "2", "boxes", "")
		//	check the result
		if want != got.Error() {
			t.Errorf("failed creating a new set of points: expected error: %s result: %s", want, got)
		}
	})

	t.Run(">>> newSetPoints2D: invalid style", func(t *testing.T) {

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

		want := `invalid style: circles`

		_, got := newSetPoints2D(tmpDataFile.Name(), "1", "2", "circles", "")
		//	check the result
		if want != got.Error() {
			t.Errorf("failed creating a new set of points: expected error: %s result: %s", want, got)
		}
	})

	t.Run(">>> newSetPoints2D: valid scenario (default title)", func(t *testing.T) {

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

		expectedPoints := 2
		expectedStyle := "boxes"
		expectedTitle := tmpDataFile.Name() + " u 1:2"

		setOfPoints, err := newSetPoints2D(tmpDataFile.Name(), "1", "2", expectedStyle, "")
		if err != nil {
			t.Errorf("fail creating a new set of points: %s", err.Error())
			return
		}

		want := expectedPoints
		got := len(setOfPoints.Point)
		//	check the result
		if want != got {
			t.Errorf("failed creating a new set of points: expected: %d points result: %d", want, got)
			return
		}

		want = int(Style[expectedStyle])
		got = int(setOfPoints.Style)
		//	check the result
		if want != got {
			t.Errorf("failed creating a new set of points: expected: %d (%s) result: %d", want, expectedStyle, got)
		}

		wantString := expectedTitle
		gotString := setOfPoints.Title
		//	check the result
		if wantString != gotString {
			t.Errorf("failed creating a new set of points: expected: %s result: %s", wantString, gotString)
		}
	})

	t.Run(">>> newSetPoints2D: valid scenario", func(t *testing.T) {

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

		expectedPoints := 2
		expectedStyle := "boxes"
		expectedTitle := "annual inflation"

		setOfPoints, err := newSetPoints2D(tmpDataFile.Name(), "1", "2", expectedStyle, expectedTitle)
		if err != nil {
			t.Errorf("fail creating a new set of points: %s", err.Error())
			return
		}

		want := expectedPoints
		got := len(setOfPoints.Point)
		//	check the result
		if want != got {
			t.Errorf("failed creating a new set of points: expected: %d points result: %d", want, got)
			return
		}

		want = int(Style[expectedStyle])
		got = int(setOfPoints.Style)
		//	check the result
		if want != got {
			t.Errorf("failed creating a new set of points: expected: %d (%s) result: %d", want, expectedStyle, got)
		}

		wantString := expectedTitle
		gotString := setOfPoints.Title
		//	check the result
		if wantString != gotString {
			t.Errorf("failed creating a new set of points: expected: %s result: %s", wantString, gotString)
		}
	})
}
