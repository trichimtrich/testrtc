package main

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{} // use default options
var clientPool map[string]*websocket.Conn


func runServer(host string) {
	log.Println("Run as 'SERVER' at:", host)
	log.Printf("Access directory http://%s/file/", host)

	clientPool = make(map[string]*websocket.Conn)

	http.HandleFunc("/", serverHandler)
	http.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir("./file"))))

	log.Fatal(http.ListenAndServe(host, nil))
}


func serverHandler(w http.ResponseWriter, r *http.Request) {
	// make sure cross origin work!
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("[!] Cannot Upgrade:", err)
		return
	}

	// generate unique ID for each client
	var myID string
	for {
		myID = generateID()
		if _, ok := clientPool[myID]; !ok {
			break
		}
	}

	// add connection ID to pool
	clientPool[myID] = conn

	defer func() {
		// remove connection ID in pool
		if myID != "" {
			if _, ok := clientPool[myID]; ok {
				delete(clientPool, myID)
			}
		}
		log.Printf("<%s> WebSocket closed\n", myID)
		conn.Close()
	}()


	log.Printf("<%s> New WebSocket created\n", myID)

	req := WSPacket{}

	// main loop
	Loop:
	for {
		err = recvMessage(conn, &req)
		if err != nil {
			log.Printf("<%s> Cannot recvMessage: %s\n", myID, err)
			break
		}

		switch (req.ID) {
		case "ping":
			err = sendMessage(conn, WSPacket{ID: "pong"})
			if err != nil {
				log.Printf("<%s> [!] Cannot sendMessage[ping]: %s\n", myID, err)
				break Loop
			}

		case "hello":
			log.Printf("<%s> Got hello packet\n", myID)
			err = sendMessage(conn, WSPacket{ID: "hello", Data: myID})
			if err != nil {
				log.Printf("<%s> [!] Cannot sendMessage[hello]: %s\n", myID, err)
				break Loop
			}

		case "mail":
			if conn2, ok := clientPool[req.ClientID]; ok {
				log.Printf("<%s> Send mail to <%s>\n", myID, req.ClientID)
				err = sendMessage(conn2, WSPacket{
					ID: "mail",
					Data: req.Data,
					ClientID: myID,
				})
				if err != nil {
					log.Printf("<%s> [!] Cannot send mail to [%s]: %s\n", myID, req.ClientID, err)
					err = sendMessage(conn, WSPacket{ID: "error", Data: "error while sending mail to partner"})
					if err != nil {
						log.Printf("<%s> [!] Cannot sendMessage[error]: %s\n", myID, err)
						break Loop
					}
				}
			} else {
				err = sendMessage(conn, WSPacket{ID: "error", Data: "invalid partner ID"})
				if err != nil {
					log.Printf("<%s> [!] Cannot sendMessage[error]: %s\n", myID, err)
					break Loop
				}
			}

		}
	}
}
