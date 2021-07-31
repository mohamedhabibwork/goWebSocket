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
	var AppName = request.FormValue("app_name")
	var channel = request.FormValue("channel")

	if channel == "" {
		return
	}

	if AppName == "" {
		AppName = "Other"
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
		defer conn.Close()
		return
	}

	uuid, _ := uuid.NewRandom()
	connId := uuid.String()
	addConnToHub(AppName, channel, connId, conn)
	go readData(AppName, channel, connId, conn)
	for {
		select {
		case message := <-MessageChan:
			for _, connWrite := range hub[message.ToClient] {
				err := connWrite.WriteJSON(message)
				if err != nil {
					removeConnToHub(AppName, channel, connId)
					defer conn.Close()
				}
			}
		}
	}
}

func addConnToHub(AppName string, channel string, connId string, conn *websocket.Conn) {
	_, ok := hub[channel]
	if ok {
		hub[channel][connId] = conn
	} else {
		connHub := make(map[string]*websocket.Conn)
		hub[channel] = connHub
		hub[channel][connId] = conn
	}
}

func removeConnToHub(AppName string, channel string, connId string) {
	_, ok := hub[channel][connId]

	if ok {
		delete(hub[channel], connId)
	}
}

func readData(AppName string, channel string, connId string, conn *websocket.Conn) {
	for {
		_, _, err := conn.ReadMessage()

		if err != nil {
			removeConnToHub(AppName, channel, connId)
			defer conn.Close()
			return
		}
	}
}

type Message struct {
	AppName    string `json:"app_name" `
	ToClient   string `json:"to_client"`
	FormClient string `json:"form_client"`
	Data       string `json:"data"`
	Type       string `json:"type"`
}
