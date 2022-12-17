////////////////////////////////////////////////////////////////////////////////
//	plotEndpoint.go  -  Sep-27-2022  -  aldebap
//
//	Plot endpoint controller
////////////////////////////////////////////////////////////////////////////////

package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"

	plot "github.com/aldebap/go-plot/plot"
)

//	plot request
type plot2DRequest struct {
	X_label string           `json:"x_label"`
	Y_label string           `json:"y_label"`
	Plot    []plotDefinition `json:"plot"`
	Width   int64            `json:"width"`
	Height  int64            `json:"height"`
}

type plotDefinition struct {
	Title        string           `json:"title"`
	DataSet      dataSetPlot      `json:"data_set"`
	MathFunction mathFunctionPlot `json:"math_function"`
}

type dataSetPlot struct {
	Points []plotPoint `json:"points"`
	Style  string      `json:"style"`
}

type plotPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type mathFunctionPlot struct {
	Min_x    float64 `json:"min_x"`
	Max_x    float64 `json:"max_x"`
	Function string  `json:"function"`
}

//	PlotHandler handle the HTTP request to generate a Go-Plot graphic
func PlotHandler(httpResponse http.ResponseWriter, httpRequest *http.Request, terminal uint8) {

	//	check for "json" content type
	contentType := httpRequest.Header.Get("Content-type")
	if contentType != "application/json" {
		httpResponse.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	//	check for non empty content length
	contentLength := httpRequest.Header.Get("Content-Length")
	if len(contentLength) == 0 {
		httpResponse.WriteHeader(http.StatusLengthRequired)
		return
	}

	if httpRequest.ContentLength == 0 {
		httpResponse.WriteHeader(http.StatusBadRequest)
		return
	}

	//	fetch request payload
	var requestData plot2DRequest

	err := json.NewDecoder(httpRequest.Body).Decode(&requestData)
	if nil != err {
		httpResponse.WriteHeader(http.StatusBadRequest)
		return
	}

	//	create a plot request from the request payload
	plotRequest := &plot.Plot_2D{
		X_label:    requestData.X_label,
		Y_label:    requestData.Y_label,
		Set_points: make([]plot.Set_points_2d, 0),
		Function:   make([]plot.Function_2d, 0),
		Width:      requestData.Width,
		Height:     requestData.Height,
		Terminal:   terminal,
	}

	for _, plotDefinition := range requestData.Plot {

		if len(plotDefinition.DataSet.Points) == 0 && len(plotDefinition.MathFunction.Function) == 0 {
			httpResponse.WriteHeader(http.StatusBadRequest)
			httpResponse.Write([]byte(fmt.Sprintf(`{ "error": "each plot must contain at least one function or one data set" }`)))
			return
		}

		if len(plotDefinition.DataSet.Points) > 0 && len(plotDefinition.MathFunction.Function) > 0 {
			httpResponse.WriteHeader(http.StatusBadRequest)
			httpResponse.Write([]byte(fmt.Sprintf(`{ "error": "each plot must be either function or data set" }`)))
			return
		}

		//	add a new set of points
		if len(plotDefinition.DataSet.Points) > 0 {

			//	attempt to convert the style string to an int constant
			var num_style uint8
			var found bool

			if len(plotDefinition.DataSet.Style) == 0 {
				num_style = plot.POINTS
			} else {
				num_style, found = plot.Style[plotDefinition.DataSet.Style]
				if !found {
					httpResponse.WriteHeader(http.StatusBadRequest)
					httpResponse.Write([]byte(fmt.Sprintf(`{ "error": "invalid style: %s" }`, plotDefinition.DataSet.Style)))
					return
				}
			}

			set_Points := plot.Set_points_2d{}

			//	set a default title when necessary
			title := plotDefinition.Title

			if len(title) == 0 {
				title = fmt.Sprintf("data set #%d", len(plotDefinition.DataSet.Points)+1)
			}

			set_Points.Title = title
			set_Points.Style = num_style

			//	add the points
			set_Points.Point = make([]plot.Point_2d, len(plotDefinition.DataSet.Points))

			for i, point := range plotDefinition.DataSet.Points {
				set_Points.Point[i].X = point.X
				set_Points.Point[i].Y = point.Y
			}

			plotRequest.Set_points = append(plotRequest.Set_points, set_Points)
		}

		//	add a new function
		if len(plotDefinition.MathFunction.Function) > 0 {

			function := plot.Function_2d{}

			//	set a default title when necessary
			title := plotDefinition.Title

			if len(title) == 0 {
				title = plotDefinition.MathFunction.Function
			}

			function.Title = title
			function.Style = plot.DOTS
			function.Function = plotDefinition.MathFunction.Function
			function.Min_x = plotDefinition.MathFunction.Min_x
			function.Max_x = plotDefinition.MathFunction.Max_x

			plotRequest.Function = append(plotRequest.Function, function)
		}
	}

	//	generate the SVG graphics as a response to HTTP request
	err = plotRequest.GeneratePlot(bufio.NewWriter(httpResponse))
	if err != nil {
		httpResponse.WriteHeader(http.StatusInternalServerError)
		httpResponse.Write([]byte(fmt.Sprintf(`{ "error": "%s" }`, err)))
		return
	}

	//	based on terminal, add the appropriate response content type
	switch terminal {
	case plot.TERMINAL_CANVAS:
		httpResponse.Header().Add("Content-Type", "text/javascript")

	case plot.TERMINAL_SVG:
		httpResponse.Header().Add("Content-Type", "image/svg+xml")

	case plot.TERMINAL_GIF:
		httpResponse.Header().Add("Content-Type", "image/gif")

	case plot.TERMINAL_JPEG:
		httpResponse.Header().Add("Content-Type", "image/jpeg")

	case plot.TERMINAL_PNG:
		httpResponse.Header().Add("Content-Type", "image/png")
	}

	httpResponse.WriteHeader(http.StatusOK)
}
