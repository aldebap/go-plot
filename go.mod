module github.com/aldebap/go-plot

go 1.17

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	golang.org/x/image v0.1.0 // indirect
)

require github.com/aldebap/go-plot/api v0.0.0-unpublished

require github.com/aldebap/go-plot/api/controller v0.0.0-unpublished

require github.com/aldebap/go-plot/expression v0.0.0-unpublished

require github.com/aldebap/go-plot/plot v0.0.0-unpublished

replace github.com/aldebap/go-plot/api v0.0.0-unpublished => ./api

replace github.com/aldebap/go-plot/api/controller v0.0.0-unpublished => ./api/controller

replace github.com/aldebap/go-plot/expression v0.0.0-unpublished => ./expression

replace github.com/aldebap/go-plot/plot v0.0.0-unpublished => ./plot
