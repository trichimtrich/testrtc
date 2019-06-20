#!/bin/bash

if [ -e main.go ]; then
    WORKDIR="./"
elif [ -e ../main.go ]; then
    WORKDIR="../"
else
    echo "Please change directory to bin or root of the repo"
    exit 1
fi

CURDIR=$(pwd)
cd $WORKDIR
go run . -host 0.0.0.0:5000
cd $CURDIR