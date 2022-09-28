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
)

//	main entry point for Go-Plot Rest API
func main() {

	var servicePort int

	//	CLI arguments
	flag.IntVar(&servicePort, "port", 8080, "port to listen for connections")

	flag.Parse()

	//	start the Web Server
	httpRouter := mux.NewRouter()

	httpRouter.HandleFunc("/plot/svg", controller.PlotSVG).Methods(http.MethodPost)

	http.Handle("/", httpRouter)

	//start and listen to requests
	fmt.Printf("Listening port %d\n", servicePort)

	log.Panic(http.ListenAndServe(fmt.Sprintf(":%d", servicePort), httpRouter))
}
