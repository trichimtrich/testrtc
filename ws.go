package main

import (
	"math/rand"
	"strconv"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"strings"
)

type WSPacket struct {
	ID   		string `json:"id"`
	Data	 	string `json:"data"`
	ClientID	string `json:"cid"`
}

type MailPacket struct {
	ID   		string `json:"id"`
	Data	 	string `json:"data"`
}


func sendMessage(conn *websocket.Conn, resp WSPacket) error {
	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Println("[!] Cannot json.Marshal in sendMessage")
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, respBytes)
	if err != nil {
		log.Println("[!] Cannot WriteMessage in sendMessage")
		return err
	}

	return nil
}


func recvMessage(conn *websocket.Conn, req *WSPacket) error {
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("[!] Cannot ReadMessage in recvMessage")
		return err
	}

	req.ID = ""
	err = json.Unmarshal(message, &req)
	if err != nil {
		log.Println("[!] Cannot json.Unmarshal in recvMessage")
		return err
	}

	return nil
}


func sendMail(conn *websocket.Conn, partnerID string, mailObj MailPacket) error {
	respBytes, err := json.Marshal(mailObj)
	if err != nil {
		log.Println("[!] Cannot json.Marshal in sendMail")
		return err
	}

	err = sendMessage(conn, WSPacket{ID: "mail", Data: string(respBytes), ClientID: partnerID})
	if err != nil {
		log.Println("[!] Cannot sendMessage")
		return err
	}

	return nil
}

func parseMail(reqData string, mailObj *MailPacket) error {
	err := json.Unmarshal([]byte(reqData), mailObj)
	if err != nil {
		log.Println("[!] Cannot json.Unmarshal in parseMail")
		return err
	}

	return nil
}



func generateID() string {
	// no need secure :3
	return strconv.Itoa(rand.Intn(100))
}


func sanitizeWsURL(host string) string {
	// http://localhost:5000/xyz -> localhost:5000/xyz
	idx := strings.Index(host, "://")
	if idx != -1 {
		host = host[idx + 3:]
	}

	// localhost:5000/xyz -> localhost:5000
	if strings.Index(host, "/") != -1 {
		host = strings.Split(host, "/")[0]
	}

	// localhost:5000 -> http://localhost:5000/ws
	host = "ws://" + host

	return host
}