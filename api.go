package main

import (
	"encoding/json"
	"net/http"
)

func sendMessage(response http.ResponseWriter, request *http.Request) {
	var MessageRequest Message

	request.Body = http.MaxBytesReader(response, request.Body, 1048576)

	dec := json.NewDecoder(request.Body)

	err := dec.Decode(&MessageRequest)

	if err != nil {
		print("error decoding body")
		var ResponseError = struct {
			Status bool   `json:"status"`
			Error  string `json:"error"`
		}{}
		ResponseError.Status = false
		ResponseError.Error = "not matched struct Message"
		messageError, _ := json.Marshal(ResponseError)
		_, _ = response.Write([]byte(messageError))
		return
	}

	if MessageRequest.AppName == "" {
		MessageRequest.AppName = "Other"
	}
	MessageChan <- MessageRequest
	var Response = struct {
		Status bool `json:"status"`
	}{}
	Response.Status = true
	message, err := json.Marshal(Response)

	_, err = response.Write([]byte(message))
	if err != nil {
		return
	}
}
