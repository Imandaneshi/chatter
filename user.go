package main

import "github.com/gorilla/websocket"

type User struct {

	room *Room

	// The websocket connection
	connection *websocket.Conn

	messages chan []byte

}

//Push sends new messages to user
func (user *User) Push() {

}

//Pull receives new messages and pushes it to other room members
func (user *User) Pull() {

}