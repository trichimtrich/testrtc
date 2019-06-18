From golang:1.12

RUN mkdir -p /testrtc
COPY . /testrtc/
WORKDIR /testrtc

# Install server dependencies
RUN apt-get update
RUN go get "github.com/pion/webrtc"
RUN go get "github.com/gorilla/websocket"

CMD ["go", "run", "."]