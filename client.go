package main

import (
	"time"
	"log"
	"github.com/gorilla/websocket"
	"os"
	"os/signal"
)

func runClient(host string, partnerID string) {
	// to Ctrl C
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	host = sanitizeWsURL(host)

	log.Println("Run as 'CLIENT' to:", host)
	conn, _, err := websocket.DefaultDialer.Dial(host, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer func() {
		conn.Close()
	}()

	done := make(chan struct{})

	// main Loop
	go clientLoop(conn, done, partnerID)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// send hello packet to get myID
	err = sendMessage(conn, WSPacket{ID: "hello"})
	if err != nil {
		log.Println("[!] Cannot sendMessage[hello]")
		return
	}

	// signal loop
	for {
		select {
		case <-done:
			log.Println("Done from server/client fault")
			return

		case <-ticker.C:
			err = sendMessage(conn, WSPacket{ID: "ping"})
			if err != nil {
				log.Println("[!] Cannot sendMessage[ping]:", err)
				return
			}

		case <-interrupt:
			log.Println("Interupt")
			return
		}
	}
}


func clientLoop(conn *websocket.Conn, done chan struct{}, partnerID string) {
	defer close(done)

	log.Println("Enter client loop")
	myID := ""
	req := WSPacket{}

	// if partnerID != "" -> create WebRTC first?

	for {
		err := recvMessage(conn, &req)
		if err != nil {
			log.Println("[!] Cannot recvMessage:", err)
			break
		}

		switch (req.ID) {
		case "error":
			log.Println("[!] Error from server:", req.Data)

		case "pong":
			break

		case "hello":
			myID = req.Data
			log.Println("Got hello packet with id:", myID)
		
		case "mail":
			log.Printf("Got mail from <%s>: %s\n", req.ClientID, req.Data)

		case "offer_sdp":

		}
	}
}