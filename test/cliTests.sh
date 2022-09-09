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
echo "[scenario #1] create a plot from a data file"

cat > ${DATA_FILE} <<DATA_CONTENT
Minutes Runners
135 1
140 2
145 4
150 7
155 11
160 13
165 35
170 29
DATA_CONTENT

cat > ${PLOT_FILE} <<PLOT_CONTENT
plot "${DATA_FILE}" using 1:2
PLOT_CONTENT

../bin/go-plot ${PLOT_FILE}

#   clean up temporary files
rm -rf ${DATA_FILE} ${PLOT_FILE} > /dev/null
