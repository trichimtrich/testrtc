# Use ubuntu for other clients in future
From ubuntu:bionic

RUN apt-get update
RUN apt-get install -y wget git curl python-software-properties

# Install golang
WORKDIR /tmp
RUN wget https://dl.google.com/go/go1.12.6.linux-amd64.tar.gz
RUN tar -xvf go1.12.6.linux-amd64.tar.gz
RUN mv go /usr/local

ENV GOROOT /usr/local/go
ENV GOPATH /root/gopath
ENV PATH "$GOPATH/bin:$GOROOT/bin:$PATH"

# Install go dependencies
RUN go get "github.com/pion/webrtc"
RUN go get "github.com/gorilla/websocket"

# Install nodejs
RUN curl -sL https://deb.nodesource.com/setup_10.x | sudo -E bash -
RUN apt-get install -y nodejs

# Install puppeteer
RUN npm install -g puppeteer

# our volume
RUN mkdir -p /testrtc
WORKDIR /testrtc

