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
	X_label string          `json:"x_label"`
	Y_label string          `json:"y_label"`
	Plot    []plotSetPoints `json:"plot"`
}

type plotSetPoints struct {
	Title  string      `json:"title"`
	Style  string      `json:"style"`
	Points []plotPoint `json:"points"`
}

type plotPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
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
		Set_points: make([]plot.Set_points_2d, len(requestData.Plot)),
		Terminal:   terminal,
	}

	for i, setPoints := range requestData.Plot {

		//	attempt to convert the style string to an int constant
		var num_style uint8
		var found bool

		num_style, found = plot.Style[setPoints.Style]
		if !found {
			httpResponse.WriteHeader(http.StatusBadRequest)
			httpResponse.Write([]byte(fmt.Sprintf(`{ "error": "invalid style: %s" }`, setPoints.Style)))
			return
		}

		//	set a default title when necessary
		title := setPoints.Title

		if len(title) == 0 {
			title = fmt.Sprintf("data set #%d", i+1)
		}

		plotRequest.Set_points[i].Title = title
		plotRequest.Set_points[i].Style = num_style

		//	add the points
		plotRequest.Set_points[i].Point = make([]plot.Point_2d, len(setPoints.Points))

		for j, point := range setPoints.Points {
			plotRequest.Set_points[i].Point[j].X = point.X
			plotRequest.Set_points[i].Point[j].Y = point.Y
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
	}

	httpResponse.WriteHeader(http.StatusCreated)
}
