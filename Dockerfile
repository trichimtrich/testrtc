# Use ubuntu for other clients in future
From ubuntu:bionic

RUN apt-get update

# Install golang
RUN apt-get install -y wget git
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
RUN apt-get install -y curl software-properties-common
RUN curl -sL https://deb.nodesource.com/setup_10.x | bash -
RUN apt-get install -y nodejs

# Install puppeteer
RUN npm install puppeteer
RUN apt-get install -y gconf-service libasound2 libatk1.0-0 libatk-bridge2.0-0 libc6 libcairo2 libcups2 libdbus-1-3 libexpat1 libfontconfig1 libgcc1 libgconf-2-4 libgdk-pixbuf2.0-0 libglib2.0-0 libgtk-3-0 libnspr4 libpango-1.0-0 libpangocairo-1.0-0 libstdc++6 libx11-6 libx11-xcb1 libxcb1 libxcomposite1 libxcursor1 libxdamage1 libxext6 libxfixes3 libxi6 libxrandr2 libxrender1 libxss1 libxtst6 ca-certificates fonts-liberation libappindicator1 libnss3 lsb-release xdg-utils

# our volume
RUN mkdir -p /testrtc
WORKDIR /testrtc

