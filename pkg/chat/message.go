package chat

type Message struct {
	Content  string `json:"content"`
	SenderId string `json:"sender_id"`
	RoomId   string `json:"room_id"`
}
