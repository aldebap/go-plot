#!  /usr/bin/ksh

#   color constants
export  RED='\033[0;31m'
export  GREEN='\033[0;32m'
export  LIGHTGRAY='\033[0;37m'
export  NOCOLOR='\033[0m'

#   set environment
export  VERBOSE='false'

#   function to execute the "run" target action
function runPackageTarget {

    local   DOCKER_FLAGS=''

    echo -e "[build] ${TARGET}: ${GREEN}run the package from the Docker image${NOCOLOR}"

    docker run -p ${PORT_MAPPING_FLAGS} $( echo ${PROJECT_TARGET} | tr [:upper:] [:lower:] ) ${DOCKER_FLAGS}
}

TARGET=run
PROJECT_TARGET=webGoPlot
PORT_MAPPING_FLAGS='6502:8080'

runPackageTarget
