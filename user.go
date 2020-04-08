package main

import "github.com/gorilla/websocket"

var users map[string]*User

type User struct {

	// the room user is in
	room string

	// The websocket connection
	connection *websocket.Conn

	messages chan []byte

	username string

	// unique user id
	id string
}

//Push sends new messages to user
func (user *User) Push() {

}

//Pull receives new messages and pushes it to other room members
func (user *User) Pull() {

}
