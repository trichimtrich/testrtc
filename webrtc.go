package main

import (
	"log"
	"encoding/base64"
	"github.com/pion/webrtc"
	"encoding/json"

)

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


func decodeIceCandidate(data string, ice *webrtc.ICECandidateInit) error {
	err := json.Unmarshal([]byte(data), ice)
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