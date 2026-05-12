package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{conn: conn}
}

func (c *Client) ReadLoop() {
	defer c.conn.Close()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("client disconnected:", err)
			break
		}

		log.Printf("message received: %s\n", string(message))
	}
}
