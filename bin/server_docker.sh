#!/bin/bash

if [ -e main.go ]; then
    WORKDIR="./"
elif [ -e ../main.go ]; then
    WORKDIR="../"
else
    echo "[!] Please change directory to bin or root of the repo"
    exit 1
fi

if [ $(id -u) != 0 ]; then
    echo "[!] Please run as root"
    exit 1
fi

docker build --no-cache -t testrtc-img $WORKDIR
docker run -p 5000:5000 -it testrtc-img go run .
