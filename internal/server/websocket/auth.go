package websocket

func CheckRoomExists(id string) bool {
	return GetHub().Rooms[id] != nil
}

func CheckUsernameAvailable(username string, roomID string) bool {
	for s := range GetHub().Rooms[roomID] {
		if s.UserID == username {
			return false
		}
	}
	return true
}
