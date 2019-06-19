package main

import (
	"time"
	"log"
	"github.com/gorilla/websocket"
	"os"
	"os/signal"
	"github.com/pion/webrtc"
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
	log.Println("[WEBRTC] Start new webrtc client")

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Printf("[WEBRTC] >>> ICE Connection State has changed: [%s]\n", connectionState.String())
	})

	// Get Ice Candidate
	peerConnection.OnICECandidate(func(iceCandidate *webrtc.ICECandidate) {
		log.Println("[WEBRTC] Have candidate:", iceCandidate)
		// var x webrtc.ICECandidateInit
	})

	// if there is a partnerID , we should be in createOffer role
	if partnerID != "" {
		offer, err := peerConnection.CreateOffer(nil)
		if err != nil {
			log.Println("[!] Cannot CreateOffer:", err)
			return
		}
		log.Println("[WEBRTC] Create Offer")

		err = peerConnection.SetLocalDescription(offer)
		if err != nil {
			log.Println("[!] Cannot SetLocalDescription:", err)
			return
		}
		log.Println("[WEBRTC] Set Local Description")

		localSession, err := encodeSDP(offer)
		if err != nil {
			log.Println("[!] Cannot encodeSDP:", err)
			return
		}
		log.Println("[WEBRTC} >>> Local SDP:", localSession)

		err = sendMail(conn, partnerID, MailPacket{ID: "offer_sdp", Data: localSession})
		if err != nil {
			log.Println("[!] Cannot sendMail:", err)
			return
		}
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
			log.Println("[WS] Got hello packet with id:", myID)
		
		case "mail":
			log.Printf("[WS] Got mail from <%s>\n", req.ClientID)
			var mailObj MailPacket
			err = parseMail(req.Data, &mailObj)
			if err != nil {
				log.Println("[!] Cannot parseMail:", err)
				break Loop
			}

			switch (mailObj.ID) {
			case "test":
				log.Println("[WS] >>> Mail test:", mailObj.Data)

			case "ice":
				log.Println("[WS] >>> Mail IceCandidate:", mailObj.Data)
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
				log.Println("[WEBRTC] AddIceCandidate:", ice.Candidate)

			case "answer_sdp":
				log.Println("[WS] >>> Got answer_sdp from other:", mailObj.Data)

				answer := webrtc.SessionDescription{}
				err = decodeSDP(mailObj.Data, &answer)
				if err != nil {
					log.Println("[!] Cannot decodeSDP:", err)
					break Loop
				}

				// Set the remote SessionDescription
				err = peerConnection.SetRemoteDescription(answer)
				if err != nil {
					log.Println("[!] Cannot SetRemoteDescription:", err)
					break Loop
				}
				log.Println("[WEBRTC] Set remote description")

				// do nothing, suppose the state is connected


			case "offer_sdp":
				log.Println("[WS] >>> Got offer_sdp from other:", mailObj.Data)

				offer := webrtc.SessionDescription{}
				err = decodeSDP(mailObj.Data, &offer)
				if err != nil {
					log.Println("[!] Cannot decodeSDP:", err)
					break Loop
				}

				// Set the remote SessionDescription
				err = peerConnection.SetRemoteDescription(offer)
				if err != nil {
					log.Println("[!] Cannot SetRemoteDescription:", err)
					break Loop
				}
				log.Println("[WEBRTC] Set remote description")

				// Create an answer
				answer, err := peerConnection.CreateAnswer(nil)
				if err != nil {
					log.Println("[!] Cannot CreateAnswer:", err)
					break Loop
				}
				log.Println("[WEBRTC] Create Answer")

				// Sets the LocalDescription, and starts our UDP listeners
				err = peerConnection.SetLocalDescription(answer)
				if err != nil {
					log.Println("[!] Cannot SetLocalDescription:", err)
					break Loop
				}
				log.Println("[WEBRTC] Set local description")

				localSession, err := encodeSDP(answer)
				if err != nil {
					log.Println("[!] Cannot encodeSDP:", err)
					return
				}
				log.Println("[WEBRTC] >>> Local SDP:", localSession)

				err = sendMail(conn, req.ClientID, MailPacket{ID: "answer_sdp", Data: localSession})
				if err != nil {
					log.Println("[!] Cannot sendMail:", err)
					break Loop
				}
				log.Println("[WS] Send answer")
			}

		}
	}
}