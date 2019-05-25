package main

import (
	"os"
	"net/http"
	"log"
	"encoding/json"
	"encoding/base64"

	"github.com/pion/webrtc"
	"github.com/gorilla/websocket"

)

type WSPacket struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}


var upgrader = websocket.Upgrader{} // use default options


func ws(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer func() {
		log.Println("Client closed")
		c.Close()
	}()

	log.Println("New websocket client")

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		req := WSPacket{}
		err = json.Unmarshal(message, &req)

		if err != nil {
			log.Println("json decode:", err)
			break
		}

		resp := WSPacket{ID: "", Data: ""}
		switch (req.ID) {
		case "sdp":
			sdpSession, err := newRTC(req.Data)
			if err != nil {
				resp.ID = "error"
				resp.Data = "error while creating RTC client"
			} else {
				resp.ID = "sdp"
				resp.Data = sdpSession
			}
			break
		case "hoho":
			log.Println("hoho")
			break
		}

		if resp.ID != "" {
			respBytes, err := json.Marshal(resp)
			if err != nil {
				log.Println("json encode:", err)
				break
			}

			err = c.WriteMessage(websocket.TextMessage, respBytes)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
}


func main() {
	host := "localhost:5000"
	if len(os.Args) > 1 {
		host = os.Args[1]
	}
	log.Println("Server is running at:", host)
	log.Printf("Access http://%s/static/rtc.html", host)
	http.HandleFunc("/ws", ws)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Fatal(http.ListenAndServe(host, nil))
}



// Encode encodes the input in base64
func Encode(obj interface{}) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

// Decode decodes the input from base64
func Decode(in string, obj interface{}) error {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		return err
	}

	return nil
}

func newRTC(remoteSession string) (string, error) {
	// Prepare the configuration
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
		return "", err
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
	})
	

	// Parse offer
	offer := webrtc.SessionDescription{}
	err = Decode(remoteSession, &offer)
	if err != nil {
		log.Println("[!] Cannot Decode:", err)
		return "", err
	}

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		log.Println("[!] Cannot SetRemoteDescription:", err)
		return "", err
	}
	log.Println("Set remote description")

	// Create an answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		log.Println("[!] Cannot CreateAnswer:", err)
		return "", err
	}

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		log.Println("[!] Cannot SetLocalDescription:", err)
		return "", err
	}
	log.Println("Set local description")

	localSession, err := Encode(answer)
	if err != nil {
		log.Println("[!] Cannot Encode:", err)
		return "", err
	}

	return localSession, nil

	// Block forever
	// select {}
}