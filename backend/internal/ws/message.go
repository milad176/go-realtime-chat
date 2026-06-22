package ws

type Message struct {
	Type     string `json:"type"`
	RoomID   string `json:"roomId,omitempty"`
	Username string `json:"username,omitempty"`
	Content  string `json:"content,omitempty"`

	Users []string `json:"users,omitempty"`
}
