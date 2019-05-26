From golang:1.12

RUN mkdir -p /cloud-game
COPY . /cloud-game/
WORKDIR /cloud-game

# Install server dependencies
RUN apt-get update

# RUN apt-get install pkg-config libvpx-dev libopus-dev libopusfile-dev -y
RUN go get "github.com/pion/webrtc"
RUN go get "github.com/gorilla/websocket"

CMD ["go", "run", "testrtc.go"]