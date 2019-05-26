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
	log.Printf("Access http://%s/file/recv.html", host)

	clientPool = make(map[string]*websocket.Conn)

	http.HandleFunc("/ws", serverHandler)
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
	clientPool[myID] = conn

	defer func() {
		// remove id map in pool
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
	for {
		err = recvMessage(conn, &req)
		if err != nil {
			log.Printf("<%s> Cannot recvMessage: %s\n", myID, err)
			break
		}

		resp := WSPacket{ID: "", Data: "", ClientID: ""}
		switch (req.ID) {
		case "ping":
			resp.ID = "pong"
			resp.Data = ""

		case "hello":
			log.Printf("<%s> Got hello packet\n", myID)
			resp.ID = "hello"
			resp.Data = myID

		case "mail":
			if conn2, ok := clientPool[req.ClientID]; ok {
				log.Printf("<%s> Send mail to <%s>\n", myID, req.ClientID)
				err = sendMessage(conn2, WSPacket{
					ID: "mail",
					Data: req.Data,
					ClientID: myID,
				})
				if err != nil {
					log.Printf("<%s> [!] Cannot send mail to [%s]\n", myID, req.ClientID)
					resp.ID = "error"
					resp.Data = "error while sending mail"
				}
			} else {
				resp.ID = "error"
				resp.Data = "invalid client ID"
			}

		}

		if resp.ID != "" {
			err = sendMessage(conn, resp)
			if err != nil {
				log.Printf("<%s> [!] Cannot sendMessage: %s\n", myID, err)
				break
			}
		}
	}
}
