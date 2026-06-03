package ws

type Message struct {
	Type    string `json:"type"`
	RoomID  string `json:"room_id,omitempty"`
	Content string `json:"content,omitempty"`
}
