package main

import (
	"net/http"
)

func main()  {

	http.Handle(
		"/websocket/messages/conn",
		MessageConn(),
	)

	http.HandleFunc(
		"/websocket/messages/send",
		sendMessage,
		)

	print("connect")

	err := http.ListenAndServe(
		":3333",
		nil,
	)

	if err != nil {
		print("error",err.Error())
		return
	}
}