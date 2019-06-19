#!/bin/bash
docker build . -t testrtc-img
docker stop testrtc-con
docker rm testrtc-con

docker run --rm --name testrtc-con -p 5000:5000 testrtc-img
