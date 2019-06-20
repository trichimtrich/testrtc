#!/bin/bash

if [ $(id -u) != 0 ]; then
    echo "[!] Please run as root"
    exit 1
fi

docker stop $(docker ps -aq)
docker rm $(docker ps -aq)