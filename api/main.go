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

	//	CLI arguments
	flag.IntVar(&servicePort, "port", 8080, "port to listen for connections")

	flag.Parse()

	//	start the Web Server
	httpRouter := mux.NewRouter()

	httpRouter.HandleFunc("/plot/api/canvas", func(httpResponse http.ResponseWriter, httpRequest *http.Request) {
		controller.PlotHandler(httpResponse, httpRequest, plot.TERMINAL_CANVAS)
	}).Methods(http.MethodPost)
	httpRouter.HandleFunc("/plot/api/svg", func(httpResponse http.ResponseWriter, httpRequest *http.Request) {
		controller.PlotHandler(httpResponse, httpRequest, plot.TERMINAL_SVG)
	}).Methods(http.MethodPost)

	http.Handle("/", httpRouter)

	//start and listen to requests
	fmt.Printf("Listening port %d\n", servicePort)

	log.Panic(http.ListenAndServe(fmt.Sprintf(":%d", servicePort), httpRouter))
}
