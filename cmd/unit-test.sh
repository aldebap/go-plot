#!  /usr/bin/ksh

#   color constants
export  RED='\033[0;31m'
export  GREEN='\033[0;32m'
export  LIGHTGRAY='\033[0;37m'
export  NOCOLOR='\033[0m'

#   set environment
export  VERBOSE='true'

#   function to execute the "unit-test" target action
function unitTestTarget {

    #local   GO_TEST_FLAGS=''
    GO_TEST_FLAGS=''

    echo -e "[build] ${TARGET}: ${GREEN}running unit tests on package ${PACKAGE_TARGET}${NOCOLOR}"

    if [ ! -z "${PACKAGE_TARGET}" ]
    then
        if [ "${VERBOSE}" == 'true' ]
        then
            GO_TEST_FLAGS='-v'
        fi

        go test ${GO_TEST_FLAGS} ${PACKAGE_TARGET}
    fi
}

TARGET=unit-test
PACKAGE_TARGET=github.com/aldebap/go-plot/api

unitTestTarget

TARGET=unit-test
PACKAGE_TARGET=github.com/aldebap/go-plot/api/controller

unitTestTarget

TARGET=unit-test
PACKAGE_TARGET=github.com/aldebap/go-plot/expression

unitTestTarget

TARGET=unit-test
PACKAGE_TARGET=github.com/aldebap/go-plot/plot

unitTestTarget
