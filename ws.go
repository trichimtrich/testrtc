package main

import (
	"math/rand"
	"strconv"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type WSPacket struct {
	ID   		string `json:"id"`
	Data	 	string `json:"data"`
	ClientID	string `json:"cid"`
}

func sendMessage(conn *websocket.Conn, resp WSPacket) error {
	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Println("[!] Cannot json.Marshal")
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, respBytes)
	if err != nil {
		log.Println("[!] Cannot WriteMessage")
		return err
	}

	return nil
}


func recvMessage(conn *websocket.Conn, req *WSPacket) error {
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("[!] Cannot ReadMessage")
		return err
	}

	req.ID = ""
	err = json.Unmarshal(message, &req)
	if err != nil {
		log.Println("[!] Cannot json.Unmarshal")
		return err
	}

	return nil
}

func generateID() string {
	// no need secure :3
	return strconv.Itoa(rand.Intn(100))
}
