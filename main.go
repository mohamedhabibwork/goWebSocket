package main

import (
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var port = os.Getenv("PORT")

	http.Handle(
		"/websocket/messages/conn",
		MessageConn(),
	)

	http.HandleFunc(
		"/websocket/messages/send",
		sendMessage,
	)

	print("connect on url http://localhost:", port)

	err = http.ListenAndServe(
		fmt.Sprintf(":%s", port),
		nil,
	)

	if err != nil {
		log.Fatal("error")
		return
	}
}
