#!/bin/bash

if [ -e main.go ]; then
    WORKDIR=$(realpath .)
elif [ -e ../main.go ]; then
    WORKDIR=$(realpath ..)
else
    echo "[!] Please change directory to bin or root of the repo"
    exit 1
fi

if [ $(id -u) != 0 ]; then
    echo "[!] Please run as root"
    exit 1
fi

docker build -t testrtc-img $WORKDIR
docker run -v $WORKDIR:/testrtc -it testrtc-img node chrome.js $@
