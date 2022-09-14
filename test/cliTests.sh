#!  /usr/bin/ksh

#   variables
export DATA_FILE=test_data_file
export PLOT_FILE=test_plot_file

#   build go-plot cli
export  CURRENT_DIR="$( pwd )"

cd ..
go build -o ./bin/go-plot
cd ${CURRENT_DIR}

#   test scenatio #1
echo "[scenario #1] create a plot from a data file (linear)"

cat > ${DATA_FILE} <<DATA_CONTENT
Minutes Runners
0 0
5 5
10 10
15 15
20 20
25 25
30 30
35 35
40 40
45 45
50 50
DATA_CONTENT

cat > ${PLOT_FILE} <<PLOT_CONTENT
plot "${DATA_FILE}" using 1:2
PLOT_CONTENT

../bin/go-plot ${PLOT_FILE}

#   test scenatio #2
echo "[scenario #2] create a plot from a data file (quadratic)"

cat > ${DATA_FILE} <<DATA_CONTENT
Minutes Runners
0 0
5 25
10 100
15 225
20 400
25 625
30 900
35 1225
40 1600
45 2025
50 2500
DATA_CONTENT

cat > ${PLOT_FILE} <<PLOT_CONTENT
plot "${DATA_FILE}" using 1:2
PLOT_CONTENT

#../bin/go-plot ${PLOT_FILE}

#   clean up temporary files
rm -rf ${DATA_FILE} ${PLOT_FILE} > /dev/null
