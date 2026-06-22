package ws

import (
	"encoding/json"
	"log"

	"github.com/milad176/go-realtime-chat/backend/internal/repository"
)

type Hub struct {
	clients map[*Client]bool
	rooms   map[string]map[*Client]bool

	register   chan *Client
	unregister chan *Client

	broadcast   chan BroadcastMessage
	messageRepo *repository.MessageRepository
	onlineUsers map[string]map[string]bool // room -> username -> true
}

func NewHub(messageRepo *repository.MessageRepository) *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan BroadcastMessage),
		messageRepo: messageRepo,
		rooms:       make(map[string]map[*Client]bool),
		onlineUsers: make(map[string]map[string]bool),
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

				if client.username != "" {
					delete(
						h.onlineUsers[client.roomID],
						client.username,
					)
				}

				h.sendOnlineUsers(client.roomID)
			}
			delete(h.clients, client)
			close(client.send)
			log.Printf("client disconnected id=%s total_clients=%d\n", client.id, len(h.clients))

		case message := <-h.broadcast:

			roomClients := h.rooms[message.RoomID]
			for client := range roomClients {
				client.send <- message.Data
			}

			log.Println("BROADCAST EVENT:", string(message.Data))
			log.Printf("broadcast room=%s members=%d", message.RoomID, len(roomClients))

			var msg Message

			err := json.Unmarshal(message.Data, &msg)
			if err != nil {
				continue
			}

			for client := range h.clients {
				if client.roomID != msg.RoomID {
					continue
				}
				client.send <- message.Data
			}
		}
	}
}

func (h *Hub) JoinRoom(client *Client, roomID string, username string) {

	// Remove from old room
	if client.roomID != "" {
		delete(h.rooms[client.roomID], client)

		if h.onlineUsers[client.roomID] != nil {
			delete(h.onlineUsers[client.roomID], username)
		}
	}

	// Create room if it doesn't exist
	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*Client]bool)
	}

	// Create online users map for room if it doesn't exist
	if h.onlineUsers[roomID] == nil {
		h.onlineUsers[roomID] = make(map[string]bool)
	}

	// Add client to new room
	h.rooms[roomID][client] = true
	// Add user to online users list
	h.onlineUsers[roomID][username] = true

	// Update client state
	client.roomID = roomID

	// Notify clients in room about new online users
	h.sendOnlineUsers(roomID)
}

func (h *Hub) sendOnlineUsers(roomID string) {
	users := []string{}

	for username := range h.onlineUsers[roomID] {
		users = append(users, username)
	}

	payload, err := json.Marshal(Message{

		Type:   "online_users",
		RoomID: roomID,
		Users:  users,
	})
	if err != nil {
		return
	}
	for client := range h.rooms[roomID] {
		client.send <- payload
	}
}
