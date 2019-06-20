# TestRTC

Use to test WebRTC connectivity based on

- Client roles
    - createOffer (send sdp)
    - createAnswer (recv sdp)

- Client implementations
    - Vanila js - [browser + cli](https://webrtc.github.io/samples/)
    - Golang - [pion library](https://github.com/pion/webrtc)

The goal is to check `connection state` on both clients 👍. **NO** data/media channel is included.

## Modules

### 📦 Signaling

WebSocket for signaling between clients.

Every client is given a unique ID. WebSocket master will maintain connection and transfer data between clients via a simple mail protocol (sender-reciever).

Implemented in Golang using [Gorilla WebSocket](https://github.com/gorilla/websocket)

### 📦 WebRTC Clients

Currently there are only 2 implementation of clients in `vanila js` and `Golang`.

Client has to connect to master via WebSocket and negotiate RTC session with other client by its unique ID.

The main goal is to check if state of both WebRTC client is `completed` 👌.

Users can add other implementation of clients or testing roles based on the signal protocol.

## Deploy

- We have `Dockerfile` and `bin` directory, make sure to check them out
- If you want to run directly in host server, follow these steps (example in `ubuntu`)

### Dependencies

👉 **Require**

1. Install `go`, and set `GOROOT`, `GOPATH`
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

2. Install these `go` dependencies
```
go get "github.com/pion/webrtc"
go get "github.com/gorilla/websocket"
```

👉 **Optional**

3. Install `nodejs`
```
curl -sL https://deb.nodesource.com/setup_10.x | sudo -E bash -
sudo apt-get install nodejs
```

4. Install `puppeteer` module
```
npm install puppeteer
```

5. Install libraries for `chrome-headless`. [Reference](https://github.com/GoogleChrome/puppeteer/blob/master/docs/troubleshooting.md)
```
sudo apt-get install gconf-service libasound2 libatk1.0-0 libatk-bridge2.0-0 libc6 libcairo2 libcups2 libdbus-1-3 libexpat1 libfontconfig1 libgcc1 libgconf-2-4 libgdk-pixbuf2.0-0 libglib2.0-0 libgtk-3-0 libnspr4 libpango-1.0-0 libpangocairo-1.0-0 libstdc++6 libx11-6 libx11-xcb1 libxcb1 libxcomposite1 libxcursor1 libxdamage1 libxext6 libxfixes3 libxi6 libxrandr2 libxrender1 libxss1 libxtst6 ca-certificates fonts-liberation libappindicator1 libnss3 lsb-release xdg-utils wget
```

6. Good to go...


### Signal server

- Run server at default setting `localhost:5000`
```
go run .
```

- Or `public interface`
```
go run . -host 0.0.0.0:5000
```

- Or `bin/server_local.sh` - `bin/server_docker.sh` (port `5000`)

- Check out `go run . --help`

### Clients

- For web client, access `file` directory for more detail
```
http://localhost:5000/file/answer.html
http://localhost:5000/file/offer.html
...
http://localhost:5000/file/offer.html?id=<other-client-id>
```

- We support `chrome-headless`
```
node chrome.js http://localhost:5000/file/answer.html
node chrome.js http://localhost:5000/file/offer.html?id=<other-client-id>
```

- For golang client, as role `createAnswer` (wait for other sending sdp)
```
go run . -server localhost:5000
```

- For golang client, as role `createOffer` (send sdp to other)
```
go run . -server localhost:5000 -partner <other-client-id>
```

- For manual webrtc in browser, check out `file/manual` directory. [Reference](http://research.edm.uhasselt.be/jori/page/Misc/QtWebRTC.html)

## Result

`not yet ... 🤫`

## License

Feel free to contribute or do whatever you want 😊
