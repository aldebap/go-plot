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
		testSet := &Set_points_2d{}

		_, _, _, _, err := testSet.getMinMax()
		if err == nil {
			t.Errorf("error expected for this test scenario")
			return
		}

		if err.Error() != want {
			t.Errorf("failed evaluatin MinMax: expected: '%s' result: '%s'", want, err.Error())
		}
	})

	t.Run(">>> getMinMax: single point in the set", func(t *testing.T) {

		want_min_x := float64(0)
		want_min_y := float64(0)
		want_max_x := float64(0)
		want_max_y := float64(0)
		testSet := &Set_points_2d{
			Point: []Point_2d{
				{X: want_min_x, Y: want_min_y},
			},
		}

		got_min_x, got_min_y, got_max_x, got_max_y, err := testSet.getMinMax()
		if err != nil {
			t.Errorf("error not expected for this test scenario")
			return
		}

		if want_min_x != got_min_x ||
			want_min_y != got_min_y ||
			want_max_x != got_max_x ||
			want_max_y != got_max_y {
			t.Errorf("failed evaluatin MinMax: expected: %f, %f, %f, %f result: %f, %f, %f, %f",
				want_min_x, want_min_y, want_max_x, want_max_y,
				got_min_x, got_min_y, got_max_x, got_max_y)
		}
	})

	t.Run(">>> getMinMax: multiple points in the set", func(t *testing.T) {

		want_min_x := float64(-10)
		want_min_y := float64(-15)
		want_max_x := float64(20)
		want_max_y := float64(25)
		testSet := &Set_points_2d{
			Point: []Point_2d{
				{X: want_min_x, Y: 0},
				{X: 0, Y: want_min_y},
				{X: want_max_x, Y: 0},
				{X: 0, Y: want_max_y},
			},
		}

		got_min_x, got_min_y, got_max_x, got_max_y, err := testSet.getMinMax()
		if err != nil {
			t.Errorf("error not expected for this test scenario")
			return
		}

		if want_min_x != got_min_x ||
			want_min_y != got_min_y ||
			want_max_x != got_max_x ||
			want_max_y != got_max_y {
			t.Errorf("failed evaluatin MinMax: expected: %f, %f, %f, %f result: %f, %f, %f, %f",
				want_min_x, want_min_y, want_max_x, want_max_y,
				got_min_x, got_min_y, got_max_x, got_max_y)
		}
	})
}
