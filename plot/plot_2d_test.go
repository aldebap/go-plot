////////////////////////////////////////////////////////////////////////////////
//	plot_2d_test.go  -  Sep-16-2022  -  aldebap
//
//	Test cases for the generation of 2D Go-Plot
////////////////////////////////////////////////////////////////////////////////

package plot

import "testing"

//	TestGeneratePlot unit tests for GeneratePlot()
func TestGeneratePlot(t *testing.T) {

	t.Run(">>> GeneratePlot: empty set", func(t *testing.T) {
	})
}

//	TestGetMinMax unit tests for getMinMax()
func TestGetMinMax(t *testing.T) {

	t.Run(">>> getMinMax: empty set", func(t *testing.T) {

		want := "no points in the set"
		testSet := &set_points_2d{}

		_, _, _, _, err := testSet.getMinMax()
		if err == nil {
			t.Errorf("error expected for this test scenario")
			return
		}

		if err.Error() != want {
			t.Errorf("failed evaluatin MinMax: expected: '%s' result: '%s'", want, err.Error())
		}
	})
}
