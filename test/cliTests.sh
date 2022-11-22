#!  /usr/bin/ksh

#   variables
export DATA_FILE=test_data_file
export PLOT_FILE=test_plot_file

#   build go-plot cli
export  CURRENT_DIR="$( pwd )"

cd ..
go build -o ./bin/go-plot
cd ${CURRENT_DIR}

#   test scenatio #01
export SCENARIO="01"
export DESCRIPTION="create a plot from a data file (dots)"

echo "[scenario #${SCENARIO}] ${DESCRIPTION}"

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

cat > "${PLOT_FILE}.${SCENARIO}" <<PLOT_CONTENT
set terminal svg
set output "${PLOT_FILE}_${SCENARIO}.svg"

plot "${DATA_FILE}" using 1:2 with dots
PLOT_CONTENT

../bin/go-plot "${PLOT_FILE}.${SCENARIO}"

rm -rf "${PLOT_FILE}.${SCENARIO}" > /dev/null

#   test scenatio #02
export SCENARIO="02"
export DESCRIPTION="create a plot from a data file (boxes)"

echo "[scenario #${SCENARIO}] ${DESCRIPTION}"

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
175 35
180 45
185 35
190 25
195 22
200 15
205 9
210 11
215 12
220 21
225 28
230 40
235 50
240 46
245 83
250 60
255 64
260 80
265 67
270 82
275 75
280 78
285 92
290 75
295 84
300 83
305 62
310 57
315 66
320 52
325 53
330 49
335 37
340 31
345 25
350 32
355 14
360 12
365 14
370 10
375 13
380 8
385 9
390 7
395 5
400 7
405 5
410 6
415 9
420 8
DATA_CONTENT

cat > "${PLOT_FILE}.${SCENARIO}" <<PLOT_CONTENT
plot "${DATA_FILE}" using 1:2 with boxes

set terminal svg
set output "${PLOT_FILE}_${SCENARIO}.svg"
PLOT_CONTENT

../bin/go-plot "${PLOT_FILE}.${SCENARIO}"

rm -rf "${PLOT_FILE}.${SCENARIO}" > /dev/null

#   test scenatio #03
export SCENARIO="03"
export DESCRIPTION="create a plot from a data file (lines)"

echo "[scenario #${SCENARIO}] ${DESCRIPTION}"

cat > ${DATA_FILE} <<DATA_CONTENT
# Average PQR and XYZ stock price (in dollars per share) per calendar year
1975 49 162
1976 52 144
1977 67 140
1978 53 122
1979 67 125
1980 46 117
1981 60 116
1982 50 113
1983 66 96
1984 70 101
1985 91 93
1986 133 92
1987 127 95
1988 136 79
1989 154 78
1990 127 85
1991 147 71
1992 146 54
1993 133 51
1994 144 49
1995 158 43
DATA_CONTENT

cat > "${PLOT_FILE}.${SCENARIO}" <<PLOT_CONTENT
set terminal svg
set output "${PLOT_FILE}_${SCENARIO}.svg"

plot "${DATA_FILE}"
     using 1:2
     with lines
PLOT_CONTENT

../bin/go-plot "${PLOT_FILE}.${SCENARIO}"

rm -rf "${PLOT_FILE}.${SCENARIO}" > /dev/null

#   test scenatio #04
export SCENARIO="04"
export DESCRIPTION="create a plot from a data file (quadratic)"

echo "[scenario #${SCENARIO}] ${DESCRIPTION}"


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

cat > "${PLOT_FILE}.${SCENARIO}" <<PLOT_CONTENT
set terminal svg
set output "${PLOT_FILE}_${SCENARIO}.svg"

plot "${DATA_FILE}" using 1:2 with linespoints
     title "points for x^2"

set xlabel "natural numbers"
set ylabel "square of the number"
PLOT_CONTENT

../bin/go-plot "${PLOT_FILE}.${SCENARIO}"

rm -rf "${PLOT_FILE}.${SCENARIO}" > /dev/null

#   test scenatio #05
export SCENARIO="05"
export DESCRIPTION="create a multiple data set plot from a data file (lines and points)"

echo "[scenario #${SCENARIO}] ${DESCRIPTION}"

cat > ${DATA_FILE} <<DATA_CONTENT
# Average PQR and XYZ stock price (in dollars per share) per calendar year
1975 49 162
1976 52 144
1977 67 140
1978 53 122
1979 67 125
1980 46 117
1981 60 116
1982 50 113
1983 66 96
1984 70 101
1985 91 93
1986 133 92
1987 127 95
1988 136 79
1989 154 78
1990 127 85
1991 147 71
1992 146 54
1993 133 51
1994 144 49
1995 158 43
DATA_CONTENT

cat > "${PLOT_FILE}.${SCENARIO}" <<PLOT_CONTENT
set terminal svg
set output "${PLOT_FILE}_${SCENARIO}.svg"

plot "${DATA_FILE}" using 1:2 with lines,
     "${DATA_FILE}" using 1:3
PLOT_CONTENT

../bin/go-plot "${PLOT_FILE}.${SCENARIO}"

rm -rf "${PLOT_FILE}.${SCENARIO}" > /dev/null

#   test scenatio #06
export SCENARIO="06"
export DESCRIPTION="create a multiple data set plot from a data file (canvas)"

echo "[scenario #${SCENARIO}] ${DESCRIPTION}"

cat > ${DATA_FILE} <<DATA_CONTENT
# Distance Men   Women        Men Women
  100     9.58   10.49 2009-08-16 1988-07-16
  200    19.19   21.34 2009-08-20 1988-09-29
  400    43.18   47.60 1999-08-26 1985-10-06
  800   100.91  113.28 2012-08-09 1983-07-26
 1000   131.96  148.98 1999-09-05 1996-08-23
 1500   206.00  230.46 1998-07-14 1993-09-11
 3000   440.67  486.11 1996-09-01 1993-09-13
 5000   757.35  851.15 2005-05-31 2008-06-06
10000  1577.53 1771.78 2005-08-26 1993-09-08
21098  3503.00 3912.00 2010-03-21 2014-02-16
42195  7377.00 8262.00 2014-09-28 2003-04-13
DATA_CONTENT

cat > "${PLOT_FILE}.${SCENARIO}" <<PLOT_CONTENT
set terminal canvas
set output "${PLOT_FILE}_${SCENARIO}.js"

plot "${DATA_FILE}" using 1:2 with lines,
     "${DATA_FILE}" using 1:3 with linespoints
PLOT_CONTENT

../bin/go-plot "${PLOT_FILE}.${SCENARIO}"

rm -rf "${PLOT_FILE}.${SCENARIO}" > /dev/null

#   test scenatio #07
export SCENARIO="07"
export DESCRIPTION="create a plot from polinomial function (SVG)"

echo "[scenario #${SCENARIO}] ${DESCRIPTION}"

cat > "${PLOT_FILE}.${SCENARIO}" <<PLOT_CONTENT
set terminal svg
set output "${PLOT_FILE}_${SCENARIO}.svg"

plot [-4:+4] x * x - 3 * x + 2
PLOT_CONTENT

../bin/go-plot "${PLOT_FILE}.${SCENARIO}"

rm -rf "${PLOT_FILE}.${SCENARIO}" > /dev/null

#   test scenatio #08
export SCENARIO="08"
export DESCRIPTION="create a plot from polinomial function (SVG)"

echo "[scenario #${SCENARIO}] ${DESCRIPTION}"

cat > "${PLOT_FILE}.${SCENARIO}" <<PLOT_CONTENT
set terminal svg
set output "${PLOT_FILE}_${SCENARIO}.svg"

plot [-4:+4] x * x + 2 * x - 3
PLOT_CONTENT

../bin/go-plot "${PLOT_FILE}.${SCENARIO}"

rm -rf "${PLOT_FILE}.${SCENARIO}" > /dev/null

#   test scenatio #09
export SCENARIO="09"
export DESCRIPTION="create a plot from sine function (SVG)"

echo "[scenario #${SCENARIO}] ${DESCRIPTION}"

cat > "${PLOT_FILE}.${SCENARIO}" <<PLOT_CONTENT
set terminal svg
set output "${PLOT_FILE}_${SCENARIO}.svg"

plot [0:6.283184] sin(x)
PLOT_CONTENT

../bin/go-plot "${PLOT_FILE}.${SCENARIO}"

rm -rf "${PLOT_FILE}.${SCENARIO}" > /dev/null

#   test scenatio #10
export SCENARIO="10"
export DESCRIPTION="create a plot from complex sine function (SVG)"

echo "[scenario #${SCENARIO}] ${DESCRIPTION}"

cat > "${PLOT_FILE}.${SCENARIO}" <<PLOT_CONTENT
set terminal svg
set output "${PLOT_FILE}_${SCENARIO}.svg"

plot [-6.283184:6.283184] x
plot [-6.283184:6.283184] x * sin(4*x)
PLOT_CONTENT

../bin/go-plot "${PLOT_FILE}.${SCENARIO}"

rm -rf "${PLOT_FILE}.${SCENARIO}" > /dev/null

#   clean up temporary files
rm -rf ${DATA_FILE} > /dev/null
