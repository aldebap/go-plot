module github.com/aldebap/go-plot

go 1.17

require github.com/gorilla/mux v1.8.0 // indirect

require github.com/aldebap/go-plot/api v0.0.0-unpublished

require github.com/aldebap/go-plot/api/controller v0.0.0-unpublished

require github.com/aldebap/go-plot/plot v0.0.0-unpublished

replace github.com/aldebap/go-plot/api v0.0.0-unpublished => ./api

replace github.com/aldebap/go-plot/api/controller v0.0.0-unpublished => ./api/controller

replace github.com/aldebap/go-plot/plot v0.0.0-unpublished => ./plot
