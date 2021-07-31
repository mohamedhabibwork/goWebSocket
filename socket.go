package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
)

var hub = make(map[string]map[string]*websocket.Conn)

var MessageChan = make(chan Message)

func MessageConn() http.Handler {
	return http.HandlerFunc(MessagesSocket)
}

func MessagesSocket(response http.ResponseWriter, request *http.Request) {
	var channel = request.FormValue("channel")

	if channel == "" {
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		EnableCompression: true,
	}

	defer (func() {
		if err := recover(); err != nil {
			return
		}
	})()

	var conn, err = upgrader.Upgrade(response, request, nil)

	if err != nil {
		defer func(conn *websocket.Conn) {
			err := conn.Close()
			if err != nil {

			}
		}(conn)
		return
	}

	uuid, _ := uuid.NewRandom()
	connId := uuid.String()
	addConnToHub(channel, connId, conn)
	go readData(channel, connId, conn)
	for {
		select {
		case message := <-MessageChan:
			for _, connWrite := range hub[message.ToClient] {
				err := connWrite.WriteJSON(message)
				if err != nil {
					removeConnToHub(channel, connId)
					defer func(conn *websocket.Conn) {
						err := conn.Close()
						if err != nil {

						}
					}(conn)
				}
			}

		}
	}
}

func addConnToHub(channel string, connId string, conn *websocket.Conn) {
	_, ok := hub[channel]
	if ok {
		hub[channel][connId] = conn
	} else {
		connHub := make(map[string]*websocket.Conn)
		hub[channel] = connHub
		hub[channel][connId] = conn
	}
}

func removeConnToHub(channel string, connId string) {
	_, ok := hub[channel][connId]

	if ok {
		delete(hub[channel], connId)
	}
}

func readData(channel string, connId string, conn *websocket.Conn) {
	for {
		_, _, err := conn.ReadMessage()

		if err != nil {
			removeConnToHub(channel, connId)
			defer func(conn *websocket.Conn) {
				err := conn.Close()
				if err != nil {

				}
			}(conn)
			return
		}
	}
}

type Message struct {
	ToClient   string `json:"to_client"`
	FormClient string `json:"form_client"`
	Data       string `json:"data"`
	Type       string `json:"type"`
}
