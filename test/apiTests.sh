#!  /usr/bin/ksh

#   variables
export REQUEST_FILE=test_plot_file

#   build go-plot API
export  CURRENT_DIR="$( pwd )"

cd ..
go build -o ./bin/go-api ./api/main.go
cd ${CURRENT_DIR}

#    execute the REST API server
../bin/go-api &
if [ $? -ne 0 ]
then
	echo -e "[error] cannot start-up the REST API server"
	exit 1
fi

PID=$!

#   test scenatio #01
export SCENARIO="01"
export DESCRIPTION="create a plot from a set of points (boxes)"

echo "[scenario #${SCENARIO}] ${DESCRIPTION}"

cat > ${REQUEST_FILE} <<REQUEST_CONTENT
{
     "plot": [
          {
               "title": "boxes from a set of points",
               "style": "boxes",
               "points": [
                    { "x": 135, "y": 1 },
                    { "x": 140, "y": 2 },
                    { "x": 145, "y": 4 },
                    { "x": 150, "y": 7 },
                    { "x": 155, "y": 11 },
                    { "x": 160, "y": 13 },
                    { "x": 165, "y": 35 },
                    { "x": 170, "y": 29 },
                    { "x": 175, "y": 35 },
                    { "x": 180, "y": 45 },
                    { "x": 185, "y": 35 },
                    { "x": 190, "y": 25 },
                    { "x": 195, "y": 22 },
                    { "x": 200, "y": 15 },
                    { "x": 205, "y": 9 },
                    { "x": 210, "y": 11 },
                    { "x": 215, "y": 12 },
                    { "x": 220, "y": 21 },
                    { "x": 225, "y": 28 },
                    { "x": 230, "y": 40 },
                    { "x": 235, "y": 50 },
                    { "x": 240, "y": 46 },
                    { "x": 245, "y": 83 },
                    { "x": 250, "y": 60 },
                    { "x": 255, "y": 64 },
                    { "x": 260, "y": 80 },
                    { "x": 265, "y": 67 },
                    { "x": 270, "y": 82 },
                    { "x": 275, "y": 75 },
                    { "x": 280, "y": 78 },
                    { "x": 285, "y": 92 },
                    { "x": 290, "y": 75 },
                    { "x": 295, "y": 84 },
                    { "x": 300, "y": 83 },
                    { "x": 305, "y": 62 },
                    { "x": 310, "y": 57 },
                    { "x": 315, "y": 66 },
                    { "x": 320, "y": 52 },
                    { "x": 325, "y": 53 },
                    { "x": 330, "y": 49 },
                    { "x": 335, "y": 37 },
                    { "x": 340, "y": 31 },
                    { "x": 345, "y": 25 },
                    { "x": 350, "y": 32 },
                    { "x": 355, "y": 14 },
                    { "x": 360, "y": 12 },
                    { "x": 365, "y": 14 },
                    { "x": 370, "y": 10 },
                    { "x": 375, "y": 13 },
                    { "x": 380, "y": 8 },
                    { "x": 385, "y": 9 },
                    { "x": 390, "y": 7 },
                    { "x": 395, "y": 5 },
                    { "x": 400, "y": 7 },
                    { "x": 405, "y": 5 },
                    { "x": 410, "y": 6 },
                    { "x": 415, "y": 9 },
                    { "x": 420, "y": 8 }
               ]
          }
     ]
}
REQUEST_CONTENT

curl --request "POST" --header "Content-type: application/json" --data @${REQUEST_FILE} --output "test_plot_api_${SCENARIO}.svg" --progress-bar localhost:8080/plot/svg

rm -rf "${REQUEST_FILE}" > /dev/null

#   test scenatio #02
export SCENARIO="02"
export DESCRIPTION="create a multiple data set plot from a data file"

echo "[scenario #${SCENARIO}] ${DESCRIPTION}"

cat > ${REQUEST_FILE} <<REQUEST_CONTENT
{
     "plot": [
          {
               "title": "Men times",
               "style": "lines",
               "points": [
                    { "x":   100, "y":    9.58 },
                    { "x":   200, "y":   19.19 },
                    { "x":   400, "y":   43.18 },
                    { "x":   800, "y":  100.91 },
                    { "x":  1000, "y":  131.96 },
                    { "x":  1500, "y":  206.00 },
                    { "x":  3000, "y":  440.67 },
                    { "x":  5000, "y":  757.35 },
                    { "x": 10000, "y": 1577.53 },
                    { "x": 21098, "y": 3503.00 },
                    { "x": 42195, "y": 7377.00 }
               ]
          },
          {
               "title": "Women times",
               "style": "lines",
               "points": [
                    { "x":   100, "y":   10.49 },
                    { "x":   200, "y":   21.34 },
                    { "x":   400, "y":   47.60 },
                    { "x":   800, "y":  113.28 },
                    { "x":  1000, "y":  148.98 },
                    { "x":  1500, "y":  230.46 },
                    { "x":  3000, "y":  486.11 },
                    { "x":  5000, "y":  851.15 },
                    { "x": 10000, "y": 1771.78 },
                    { "x": 21098, "y": 3912.00 },
                    { "x": 42195, "y": 8262.00 }
               ]
          }
     ]
}
REQUEST_CONTENT

curl --request "POST" --header "Content-type: application/json" --data @${REQUEST_FILE} --output "test_plot_api_${SCENARIO}.js" --progress-bar localhost:8080/plot/canvas

rm -rf "${REQUEST_FILE}" > /dev/null

kill -9 ${PID}
