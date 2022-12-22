# syntax=docker/dockerfile:1

#   copy Web Go-Plot Application files to nginx directory
#FROM nginx
#RUN rm /etc/nginx/nginx.conf /etc/nginx/conf.d/default.conf
#COPY web /usr/share/nginx/html/plot
#COPY web/css /usr/share/nginx/html/plot/css
#COPY web/js /usr/share/nginx/html/plot/js
#COPY config/nginx/sites-available/container-go-plot /etc/nginx/sites-available

#	build the application
FROM golang:1.17-alpine AS build

WORKDIR /go-plot

COPY main.go go.mod go.sum ./
COPY api/main.go api/go.mod api/go.sum ./api/
COPY api/controller/plotEndpoint.go api/controller/go.mod ./api/controller/
COPY expression/expression.go expression/queue.go expression/stack.go expression/symbol.go expression/go.mod ./expression/
COPY plot/canvasDriver.go plot/graphicsDriver.go plot/imageDriver.go plot/svgDriver.go plot/dataFile.go plot/plot_2d.go plot/plot.go plot/plotFile.go plot/go.mod ./plot/

COPY web ./web
COPY web/css ./web/css
COPY web/js ./web/js

RUN CGO_ENABLED=0 go build -o ./bin/server ./api/main.go

#	create application image
FROM alpine:latest

WORKDIR /go-plot

COPY --from=build /go-plot/web ./web
COPY --from=build /go-plot/web/css ./web/css
COPY --from=build /go-plot/web/js ./web/js

COPY --from=build /go-plot/bin/server ./bin/

ENTRYPOINT ["/go-plot/bin/server", "-port", "8080", "-webAppDirectory", "/go-plot/web"]

EXPOSE 8000
