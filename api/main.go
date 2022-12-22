////////////////////////////////////////////////////////////////////////////////
//	main.go  -  Sep-27-2022  -  aldebap
//
//	Entry point of a Rest API for Go-Plot
////////////////////////////////////////////////////////////////////////////////

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/aldebap/go-plot/api/controller"
	"github.com/aldebap/go-plot/plot"
)

//	main entry point for Go-Plot Rest API
func main() {

	var servicePort int
	var webAppDirectory string

	//	CLI arguments
	flag.IntVar(&servicePort, "port", 8080, "port to listen for connections")
	flag.StringVar(&webAppDirectory, "webAppDirectory", "", "directory with content for webApp")

	flag.Parse()

	//	start the Web Server
	httpRouter := mux.NewRouter()

	httpRouter.HandleFunc("/plot/api/canvas", func(httpResponse http.ResponseWriter, httpRequest *http.Request) {
		controller.PlotHandler(httpResponse, httpRequest, plot.TERMINAL_CANVAS)
	}).Methods(http.MethodPost)
	httpRouter.HandleFunc("/plot/api/svg", func(httpResponse http.ResponseWriter, httpRequest *http.Request) {
		controller.PlotHandler(httpResponse, httpRequest, plot.TERMINAL_SVG)
	}).Methods(http.MethodPost)

	httpRouter.HandleFunc("/plot/api/gif", func(httpResponse http.ResponseWriter, httpRequest *http.Request) {
		controller.PlotHandler(httpResponse, httpRequest, plot.TERMINAL_GIF)
	}).Methods(http.MethodPost)
	httpRouter.HandleFunc("/plot/api/jpeg", func(httpResponse http.ResponseWriter, httpRequest *http.Request) {
		controller.PlotHandler(httpResponse, httpRequest, plot.TERMINAL_JPEG)
	}).Methods(http.MethodPost)
	httpRouter.HandleFunc("/plot/api/png", func(httpResponse http.ResponseWriter, httpRequest *http.Request) {
		controller.PlotHandler(httpResponse, httpRequest, plot.TERMINAL_PNG)
	}).Methods(http.MethodPost)

	//	if informed, add a handler for webApp content
	if len(webAppDirectory) > 0 {
		httpRouter.PathPrefix("/plot/").Handler(http.StripPrefix("/plot/", http.FileServer(http.Dir(webAppDirectory))))
	}

	http.Handle("/", httpRouter)

	//	start and listen to requests
	fmt.Printf("Listening port %d\n", servicePort)

	log.Panic(http.ListenAndServe(fmt.Sprintf(":%d", servicePort), httpRouter))
}
