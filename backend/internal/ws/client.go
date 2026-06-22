package ws

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	id       string
	roomID   string
	username string

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
			c.username = message.Username
			c.hub.JoinRoom(c, message.RoomID, message.Username)

			log.Printf("client=%s joined room=%s username=%s\n", c.id, message.RoomID, message.Username)

		case "typing":
			roomClients := c.hub.rooms[c.roomID]

			for client := range roomClients {
				if client.id == c.id {
					continue
				}

				client.send <- payload
			}

		case "chat_message":
			if c.roomID == "" {
				log.Println("client not in room")
				continue
			}

			err := c.hub.messageRepo.Create(c.roomID, message.Username, message.Content)
			if err != nil {
				log.Println("failed to save message:", err)
				continue
			}

			c.hub.broadcast <- BroadcastMessage{
				RoomID: c.roomID,
				Data:   payload,
			}
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
