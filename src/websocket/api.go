package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func sendMessage(response http.ResponseWriter, request *http.Request) {
	var MessageRequest Message

	request.Body = http.MaxBytesReader(response, request.Body, 1048576)

	dec := json.NewDecoder(request.Body)

	err := dec.Decode(&MessageRequest)

	if err != nil {
		fmt.Println("error decoding body")
		return
	}

	MessageChan <- MessageRequest
	var Response = struct {
		Status bool `json:"status"`
	}{}
	Response.Status=true
	message,err:=json.Marshal(Response)

	response.Write([]byte(message))
}
