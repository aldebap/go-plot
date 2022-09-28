module github.com/aldebap/go-plot/api

go 1.17

require github.com/aldebap/go-plot/api/controller v0.0.0-unpublished

require github.com/aldebap/go-plot/plot v0.0.0-unpublished

require github.com/gorilla/mux v1.8.0 // indirect

replace github.com/aldebap/go-plot/api/controller v0.0.0-unpublished => ./controller

replace github.com/aldebap/go-plot/plot v0.0.0-unpublished => ../plot
