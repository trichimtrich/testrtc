package main

import (
	"encoding/base64"
	"github.com/pion/webrtc"
	"encoding/json"
)

// Encode SDP: json -> base64 -> send
func encodeSDP(obj interface{}) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

// Decode SDP: recv -> base64 -> json
func decodeSDP(in string, obj interface{}) error {
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


// Decode Ice Candidate: json
func decodeIceCandidate(data string, ice *webrtc.ICECandidateInit) error {
	err := json.Unmarshal([]byte(data), ice)
	if err != nil {
		return err
	}

	return nil
}