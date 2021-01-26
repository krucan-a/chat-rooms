package websocket

func GetUsername(userID string, roomID string) (username string, ok bool) {
	hub := GetHub()
	for sub := range hub.Rooms[roomID] {
		if sub.UserID == userID {
			return sub.Username, true
		}
	}
	return "", false
}
