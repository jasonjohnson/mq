package main

import "io/ioutil"
import "net/http"

func GetQueue(session *Session) {
	session.Response.WriteHeader(http.StatusOK)
}

func CreateQueue(session *Session) {
	queue := &Queue{
		Id: session.Match.Variables["queue"],
	}

	session.Store.SaveQueue(queue)
	session.Response.WriteHeader(http.StatusCreated)
}

func DeleteQueue(session *Session) {
	queue := &Queue{
		Id: session.Match.Variables["queue"],
	}

	session.Store.DeleteQueue(queue)
	session.Response.WriteHeader(http.StatusOK)
}

func GetMessages(session *Session) {
	// ...
}

func CreateMessage(session *Session) {
	queue := &Queue{
		Id: session.Match.Variables["queue"],
	}

	body, err := ioutil.ReadAll(session.Request.Body)

	if err != nil {
		session.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	message := &Message{
		Id:      getRandomUUID(),
		Content: body,
	}

	session.Store.SaveMessage(queue, message)
	session.Response.WriteHeader(http.StatusCreated)
}

func DeleteMessage(session *Session) {
	// ...
}
