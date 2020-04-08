package main

import (
	"github.com/gorilla/websocket"
	uuid "github.com/nu7hatch/gouuid"
	log "github.com/sirupsen/logrus"
	"time"
)

var users map[string]*user

type user struct {

	// the room user is in
	room string

	// The websocket connection
	connection *websocket.Conn

	messages chan *message

	Username string

	// unique user ID
	ID string
}

//Push sends new messages to user
func (u *user) push() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		u.connection.Close()
		u.removeFromUsers()
		log.Infof("user disconnected: %s - %s", u.ID, u.Username)
	}()
	for {
		select {
		case message, ok := <-u.messages:
			u.connection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				u.connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			_ = u.connection.WriteJSON(message.toProtectedMap())

			// Add queued chat messages to the current websocket message.
			n := len(u.messages)
			for i := 0; i < n; i++ {
				msg := <-u.messages
				_ = u.connection.WriteJSON(msg.toProtectedMap())
			}
		case <-ticker.C:
			u.connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := u.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

//Pull receives new messages from user and pushes it to other room members
func (u *user) pull() {
	defer func() {
		_ = u.connection.Close()
		u.removeFromUsers()
		log.Infof("user disconnected: %s - %s", u.ID, u.Username)
	}()
	u.connection.SetReadLimit(maxMessageSize)
	u.connection.SetReadDeadline(time.Now().Add(pongWait))
	u.connection.SetPongHandler(func(string) error { u.connection.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := u.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Info("user disconnected: %v", err)
			}
			break
		}
		jsonError, msg := messageFromJson(message)

		if jsonError != nil {
			log.Infof("failed encoding message from user pull: %s", message)
			u.sendError("invalidJson")
			continue
		}

		msg.Room = u.room
		msg.UserUsername = u.Username
		msg.UserId = u.ID

		switch msg.Action {
		case "self":
			u.sendSelf()
		case "sendText":
			publish(msg)
		case "typing":
			msg.Text = ""
			publish(msg)
		case "buzz":
			msg.Text = ""
			publish(msg)
		default:
			continue
		}

	}

}

// assignId assignees a unique uuid to user
func (u *user) assignId() {
	u.ID = randId()
}

func (u *user) setUsername(username string) {
	u.Username = username
}

func (u *user) toSelfMap() map[string]interface{} {
	return map[string]interface{}{"action": "self",
		"username": u.Username, "userId": u.ID, "room": u.room}
}

// sendSelf sends user info
func (u *user) sendSelf() {
	_ = u.connection.WriteJSON(u.toSelfMap())
}

// sendError sends specified error as json
func (u *user) sendError(errorCode string) {
	err := map[string]interface{}{
		"action": "error",
		"text": errorCode,
	}
	_ = u.connection.WriteJSON(err)
}

// addToUsers adds user to users list
func (u *user) addToUsers() {
	if users == nil {
		users = map[string]*user{}
	}
	users[u.ID] = u
}

// removeFromUsers removes user from users list
func (u *user) removeFromUsers(){
	if _, ok := users[u.ID]; ok {
		delete(users, u.ID)
	}
}

// randId generate a random uuid
func randId() string {
	u, _ := uuid.NewV4()
	return u.String()
}