////////////////////////////////////////////////////////////////////////////////
//	plotEndpoint.go  -  Sep-27-2022  -  aldebap
//
//	Plot endpoint controller
////////////////////////////////////////////////////////////////////////////////

package controller

import (
	"encoding/json"
	"net/http"
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

//	PlotSVG generate a SVG graphic based on plot request
func PlotSVG(httpResponse http.ResponseWriter, httpRequest *http.Request) {

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

	//	create a plot from the request payload

	//	fill response payload
	httpResponse.Header().Add("Content-Type", "application/json")
	httpResponse.WriteHeader(http.StatusCreated)
}
