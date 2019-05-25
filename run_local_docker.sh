#!/bin/bash
docker build . -t cloud-game-local
docker stop cloud-game-local
docker rm cloud-game-local

# Overlord and worker should be run separately. Local is for demo purpose
docker run --privileged -d --name cloud-game-local -p 5000:5000 cloud-game-local bash -c "go run testrtc.go"
