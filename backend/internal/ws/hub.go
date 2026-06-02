package ws

import (
	"encoding/json"
	"log"
)

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	rooms      map[string]map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		rooms:      make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {

		select {

		case client := <-h.register:
			h.clients[client] = true
			log.Printf("client connected id=%s total_clients=%d\n", client.id, len(h.clients))

		case client := <-h.unregister:
			if client.roomID != "" {
				delete(h.rooms[client.roomID], client)
			}
			delete(h.clients, client)
			close(client.send)
			log.Printf("client disconnected id=%s total_clients=%d\n", client.id, len(h.clients))

		case message := <-h.broadcast:
			log.Println("BROADCAST EVENT:", string(message))
			log.Println("CLIENTS COUNT:", len(h.clients))

			var msg Message

			err := json.Unmarshal(message, &msg)
			if err != nil {
				continue
			}

			for client := range h.clients {
				if client.roomID != msg.RoomID {
					continue
				}
				client.send <- message
			}
		}
	}
}

func (h *Hub) JoinRoom(client *Client, roomID string) {

	// Remove from old room
	if client.roomID != "" {
		delete(h.rooms[client.roomID], client)
	}

	// Create room if it doesn't exist
	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*Client]bool)
	}

	// Add client to new room
	h.rooms[roomID][client] = true

	// Update client state
	client.roomID = roomID
}
