# TestRTC

Use to test WebRTC connectivity based on
- Client roles 
    - createOffer (send sdp)
    - createAnswer (recv sdp)

- Client implementations
    - Golang - [pion library](https://github.com/pion/webrtc)
    - Vanila js - [browser](https://webrtc.github.io/samples/)

## Modules

### Signaling

WebSocket for signaling between clients.

Every client is given a unique ID. WebSocket master will maintain connection and transfer data between clients via a simple mail protocol (sender-reciever).

Implemented in Golang using [Gorilla WebSocket](https://github.com/gorilla/websocket)

### WebRTC Client

Currently there are only 2 implementation of clients in `Golang` and `vanila js`.

Every client has to connect to master via WebSocket and negotiate RTC session with other client by its unique ID.

The main goal is to check if state of both WebRTC client is "completed".

Users can add other implementation of clients or testing roles based on the signal protocol.

## Deploy

### Signal server

- Run server at default setting `localhost:5000`
```
go run .
```

- Or
```
go run . 0.0.0.0:5000
```

- Or using `docker-compose` (port `5000`)
```
docker-compose up
```

### Clients

- For web client, access `file` directory for more detail
```
http://localhost:5000/file/
```

- For golang client, role `createAnswer`
```
go run . -server localhost:5000
```

- For golang client, role `createOffer`
```
go run . -server localhost:5000 -partnerID <other-client-id>
```

## Result

`not yet ...`

## License

Feel free to contribute or do whatever you want 😊