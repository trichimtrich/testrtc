# Use ubuntu for other clients in future
From ubuntu:bionic

RUN apt-get update
RUN apt-get install -y wget

# Install golang
WORKDIR /tmp
RUN wget https://dl.google.com/go/go1.12.6.linux-amd64.tar.gz
RUN tar -xvf go1.12.6.linux-amd64.tar.gz
RUN mv go /usr/local

RUN export GOROOT=/usr/local/go
RUN export GOPATH=$HOME/go
RUN export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# Install go dependencies
RUN go get "github.com/pion/webrtc"
RUN go get "github.com/gorilla/websocket"


# our volume
RUN mkdir -p /testrtc
WORKDIR /testrtc

