package websocket

import "log"

var hub = newHub()

func GetHub() *Hub {
	return hub
}

func newHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]map[*Subscription]bool),
		Broadcast:  make(chan *Message),
		Register:   make(chan *Subscription),
		Unregister: make(chan *Subscription),
	}
}

func (h *Hub) Run() {
	for {
		select {
		// add user to room
		case sub := <-h.Register:
			h.addUserToRoom(sub)
		// remove user from room
		case sub := <-h.Unregister:
			h.removeUserFromRoom(sub)
		// send message to everyone in the room
		case msg := <-h.Broadcast:
			h.broadcastMessageInRoom(msg)
		}
	}
}

func (h *Hub) addUserToRoom(sub *Subscription) {
	log.Printf("Adding user %s to room %s", sub.UserID, sub.RoomID)
	h.Rooms[sub.RoomID][sub] = true
}

func (h *Hub) removeUserFromRoom(sub *Subscription) {
	subscriptions := h.Rooms[sub.RoomID]
	if subscriptions != nil {
		if _, ok := subscriptions[sub]; ok {
			log.Printf("Removing user %s from room %s", sub.UserID, sub.RoomID)
			close(sub.Connection.Send)
			delete(subscriptions, sub)
			if len(subscriptions) == 0 {
				log.Printf("Deleting room %s", sub.RoomID)
				delete(h.Rooms, sub.RoomID)
			}
		}
	}
}

func (h *Hub) broadcastMessageInRoom(msg *Message) {
	subscriptions := h.Rooms[msg.RoomID]
	for sub := range subscriptions {
		select {
		case sub.Connection.Send <-msg:
		default:
			log.Printf("Removing user %s from room %s", sub.UserID, sub.RoomID)
			close(sub.Connection.Send)
			delete(subscriptions, sub)
			if len(subscriptions) == 0 {
				log.Printf("Deleting room %s", sub.RoomID)
				delete(h.Rooms, msg.RoomID)
			}
		}
	}
}
