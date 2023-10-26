package chat

type Message struct {
	Content    string          `json:"content"`
	SenderId   int             `json:"sender_id"`
	SenderName string          `json:"sender_name"`
	RoomId     string          `json:"room_id"`
	Meta       map[MetaKey]any `json:"meta"`
}

type MetaKey string

const (
	MkRoomId MetaKey = "room_id"
)

func (m *Message) WithMeta(data map[MetaKey]any) *Message {
	m.Meta = data
	return m
}

func (m *Message) AddMeta(key MetaKey, value any) {
	m.Meta[key] = value
}
