////////////////////////////////////////////////////////////////////////////////
//	dataFile_test.go  -  Sep-6-2022  -  aldebap
//
//	Test cases for parser for Go-Plot data files
////////////////////////////////////////////////////////////////////////////////

package plot

import (
	"bufio"
	"strings"
	"testing"
)

//	TestLoadDataFile unit tests for LoadDataFile()
func TestLoadDataFile(t *testing.T) {

	t.Run(">>> LoadDataFile: ignoring the header line", func(t *testing.T) {

		mockDataFile := strings.NewReader("col1 col2 col3\n")
		expectedPoints := 0
		point, err := LoadDataFile(1, 2, bufio.NewReader(mockDataFile))
		if err != nil {
			t.Errorf("fail loading data file: %s", err.Error())
			return
		}

		want := expectedPoints
		got := len(point)
		//	check the size of result
		if want != got {
			t.Errorf("failed parsing data file: size expected: %d result: %d", want, got)
		}
	})

	t.Run(">>> LoadDataFile: total lines parsed", func(t *testing.T) {

		mockDataFile := strings.NewReader("col1 col2 col3\n10 20 30\n40 50 60\n")
		expectedPoints := 2
		point, err := LoadDataFile(1, 2, bufio.NewReader(mockDataFile))
		if err != nil {
			t.Errorf("fail loading data file: %s", err.Error())
			return
		}

		want := expectedPoints
		got := len(point)
		//	check the size of result
		if want != got {
			t.Errorf("failed parsing data file: size expected: %d result: %d", want, got)
		}
	})

	t.Run(">>> LoadDataFile: single line", func(t *testing.T) {

		mockDataFile := strings.NewReader("col1 col2 col3\n10 20 30\n")
		expectedPoints := 1
		point, err := LoadDataFile(1, 2, bufio.NewReader(mockDataFile))
		if err != nil {
			t.Errorf("fail loading data file: %s", err.Error())
			return
		}

		//	check the size of result
		if len(point) != expectedPoints {
			t.Errorf("failed parsing data file: size expected: %d result: %d", expectedPoints, len(point))
			return
		}

		want := float64(10)
		got := point[0].x
		//	check the x coordinate result
		if want != got {
			t.Errorf("failed parsing data file: x coordinated expected %f result: %f", want, got)
		}

		want = float64(20)
		got = point[0].y
		//	check the y coordinate result
		if want != got {
			t.Errorf("failed parsing data file: y coordinated expected %f result: %f", want, got)
		}
	})

	t.Run(">>> LoadDataFile: double line", func(t *testing.T) {

		mockDataFile := strings.NewReader("col1 col2 col3\n10 20 30\n40 50 60\n")
		expectedPoints := 2
		point, err := LoadDataFile(2, 3, bufio.NewReader(mockDataFile))
		if err != nil {
			t.Errorf("fail loading data file: %s", err.Error())
			return
		}

		//	check the size of result
		if len(point) != expectedPoints {
			t.Errorf("failed parsing data file: size expected: %d result: %d", expectedPoints, len(point))
			return
		}

		want := float64(50)
		got := point[1].x
		//	check the x coordinate result
		if want != got {
			t.Errorf("failed parsing data file: x coordinated expected %f result: %f", want, got)
		}

		want = float64(60)
		got = point[1].y
		//	check the y coordinate result
		if want != got {
			t.Errorf("failed parsing data file: y coordinated expected %f result: %f", want, got)
		}
	})

	t.Run(">>> LoadDataFile: lines with less columns than expected", func(t *testing.T) {

		mockDataFile := strings.NewReader("col1 col2 col3\n10 20 30\n40 50 60\n")
		_, err := LoadDataFile(1, 4, bufio.NewReader(mockDataFile))
		if err == nil {
			t.Errorf("fail expected for this data file")
			return
		}

		want := `line with less columns than expected: "10 20 30"`
		got := err.Error()
		//	check the size of result
		if want != got {
			t.Errorf("failed parsing data file: error expected: '%s' result: '%s'", want, got)
		}
	})

	t.Run(">>> LoadDataFile: x column expected to be numeric", func(t *testing.T) {

		mockDataFile := strings.NewReader("col1 col2 col3\n10a 20 30\n40 50 60\n")
		_, err := LoadDataFile(1, 3, bufio.NewReader(mockDataFile))
		if err == nil {
			t.Errorf("fail expected for this data file")
			return
		}

		want := `column 1 expected to be numeric: "10a 20 30"`
		got := err.Error()
		//	check the size of result
		if want != got {
			t.Errorf("failed parsing data file: error expected: '%s' result: '%s'", want, got)
		}
	})

	t.Run(">>> LoadDataFile: y column expected to be numeric", func(t *testing.T) {

		mockDataFile := strings.NewReader("col1 col2 col3\n10 20 30c\n40 50 60\n")
		_, err := LoadDataFile(1, 3, bufio.NewReader(mockDataFile))
		if err == nil {
			t.Errorf("fail expected for this data file")
			return
		}

		want := `column 3 expected to be numeric: "10 20 30c"`
		got := err.Error()
		//	check the size of result
		if want != got {
			t.Errorf("failed parsing data file: error expected: '%s' result: '%s'", want, got)
		}
	})
}
