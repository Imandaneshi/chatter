package main

type Room struct {

	//list of registered users in this room
	users map[*User]bool

	//users to be added to this room
	add chan *User

	//users to be removed from this room
	remove chan *User

	//messages to be pushed to users
	messages chan []byte
}


func (room *Room) Start() {

	//keep looking for new messages and users
	for {
		select {
			case user := <-room.remove:
				if _,ok := room.users[user]; ok {
					delete(room.users, user)
					close(user.messages)
				}
			case user:= <-room.add:
				room.users[user] = true
		}
	}

}