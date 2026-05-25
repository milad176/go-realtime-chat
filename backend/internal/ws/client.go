package ws

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	id     string
	roomID string

	conn *websocket.Conn
	hub  *Hub
	send chan []byte
}

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		id:     uuid.New().String(),
		roomID: "",

		conn: conn,
		hub:  hub,
		send: make(chan []byte),
	}
}

func (c *Client) ReadLoop() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("client disconnected:", err)
			break
		}

		var message Message

		err = json.Unmarshal(payload, &message)

		if err != nil {
			log.Println("invalid message:", err)
			continue
		}

		log.Printf("message received with type: %s\n", string(message.Type))

		switch message.Type {

		case "join_room":
			c.roomID = message.RoomID
			log.Printf("client=%s joined room=%s\n", c.id, c.roomID)

		case "chat_message":
			c.hub.broadcast <- payload
		}
	}
}

func (c *Client) WriteLoop() {
	defer c.conn.Close()

	for {
		message, ok := <-c.send

		if !ok {
			return
		}

		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
	}
}
