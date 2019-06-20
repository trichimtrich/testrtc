# TestRTC

Use to test WebRTC connectivity based on
- Client roles
    - createOffer (send sdp)
    - createAnswer (recv sdp)

- Client implementations
    - Vanila js - [browser](https://webrtc.github.io/samples/)
    - Golang - [pion library](https://github.com/pion/webrtc)

The goal is to check `connection state` on both clients. **NO** data/media channel is included.

## Modules

### Signaling

WebSocket for signaling between clients.

Every client is given a unique ID. WebSocket master will maintain connection and transfer data between clients via a simple mail protocol (sender-reciever).

Implemented in Golang using [Gorilla WebSocket](https://github.com/gorilla/websocket)

### WebRTC Client

Currently there are only 2 implementation of clients in `vanila js` and `Golang`.

Every client has to connect to master via WebSocket and negotiate RTC session with other client by its unique ID.

The main goal is to check if state of both WebRTC client is "completed".

Users can add other implementation of clients or testing roles based on the signal protocol.

## Deploy

- We have `Dockerfile` and `bin` directory, make sure to check them out
- If you want to run directly in host server, follow these steps (example in `ubuntu`)

### Dependencies

1. Install `go`, export `GOROOT`, `GOPATH`
```
cd /tmp
wget https://dl.google.com/go/go1.12.6.linux-amd64.tar.gz
tar -xvf go1.12.6.linux-amd64.tar.gz
sudo mv go /usr/local

# ...
export GOROOT=/usr/local/go
export GOPATH=$HOME/my-go-path
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```

2. Install `go` dependencies
```
go get "github.com/pion/webrtc"
go get "github.com/gorilla/websocket"
```

3. Install `nodejs`
```
curl -sL https://deb.nodesource.com/setup_10.x | sudo -E bash -
sudo apt-get install nodejs
```

4. Install `puppeteer` module
```
sudo npm install -g puppeteer
```

5. Good to go...


### Signal server

- Run server at default setting `localhost:5000`
```
go run .
```

- Or
```
go run . -host 0.0.0.0:5000
```

- Or using `docker-compose` (port `5000`)

- Check out `go run . --help`

### Clients

- For web client, access `file` directory for more detail -> `offer.html` & `answer.html`
```
http://localhost:5000/file/
```

- For golang client, as role `createAnswer` (wait for other sending sdp)
```
go run . -server localhost:5000
```

- For golang client, as role `createOffer` (send sdp to other)
```
go run . -server localhost:5000 -partner <other-client-id>
```

- For manual webrtc in browser, check out `file\manual` directory. [Reference](http://research.edm.uhasselt.be/jori/page/Misc/QtWebRTC.html)

## Result

`not yet ... 🤫`

## License

Feel free to contribute or do whatever you want 😊
