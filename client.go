package main

import (
	"time"
	"log"
	"github.com/gorilla/websocket"
	"os"
	"os/signal"
	"github.com/pion/webrtc"
	// "encoding/json"
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

	// create WebRTC
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Println("[!] Cannot NewPeerConnection:", err)
		return
	}
	log.Println("Start new webrtc client")

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Printf("ICE Connection State has changed: [%s]\n", connectionState.String())
	})

	// Get Ice Candidate
	peerConnection.OnICECandidate(func(iceCandidate *webrtc.ICECandidate) {
		log.Println("Have candidate:", iceCandidate)
		// var x webrtc.ICECandidateInit
	})

	// if there is a partnerID , we should be in createOffer role
	if partnerID != "" {
	}

	Loop:
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
			var mailObj MailPacket
			err = parseMail(req.Data, &mailObj)
			if err != nil {
				log.Println("[!] Cannot parseMail:", err)
				break Loop
			}

			switch (mailObj.ID) {
			case "test":
				log.Printf("> Mail test: %s\n", mailObj.Data)

			case "ice":
				var ice webrtc.ICECandidateInit
				err = decodeIceCandidate(mailObj.Data, &ice)
				if err != nil {
					log.Println("[!] Cannot decodeIceCandidate:", err)
					break Loop
				}
				
				err = peerConnection.AddICECandidate(ice)
				if err != nil {
					log.Println("[!] Cannot AddIceCandidate", err)
					break Loop
				}
				log.Println("AddIceCandidate:", ice.Candidate)

			case "offer_sdp":
				offer := webrtc.SessionDescription{}
				err = Decode(mailObj.Data, &offer)
				if err != nil {
					log.Println("[!] Cannot Decode:", err)
					break Loop
				}

				// Set the remote SessionDescription
				err = peerConnection.SetRemoteDescription(offer)
				if err != nil {
					log.Println("[!] Cannot SetRemoteDescription:", err)
					break Loop
				}
				log.Println("Set remote description")

				// Create an answer
				answer, err := peerConnection.CreateAnswer(nil)
				if err != nil {
					log.Println("[!] Cannot CreateAnswer:", err)
					break Loop
				}

				// Sets the LocalDescription, and starts our UDP listeners
				err = peerConnection.SetLocalDescription(answer)
				if err != nil {
					log.Println("[!] Cannot SetLocalDescription:", err)
					break Loop
				}
				log.Println("Set local description")

				localSession, err := Encode(answer)
				if err != nil {
					log.Println("[!] Cannot Encode:", err)
					return
				}

				err = sendMail(conn, req.ClientID, MailPacket{ID: "answer_sdp", Data: localSession})
				if err != nil {
					log.Println("[!] Cannot sendMail:", err)
					break Loop
				}
				log.Println("Send answer")
			}

		}
	}
}