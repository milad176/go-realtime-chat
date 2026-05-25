package ws

type Message struct {
	Type    string `json:"type"`
	RoomID  string `json:"roomId,omitempty"`
	Content string `json:"content,omitempty"`
}
