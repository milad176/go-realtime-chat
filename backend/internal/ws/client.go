package ws

import (
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
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("client disconnected:", err)
			break
		}

		log.Printf("message received: %s\n", string(message))

		c.hub.broadcast <- message
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
