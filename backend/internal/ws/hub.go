package ws

import "log"

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {

		select {

		case client := <-h.register:
			h.clients[client] = true
			log.Println("CONNECTED CLIENTS:", len(h.clients))

		case client := <-h.unregister:
			delete(h.clients, client)
			close(client.send)

		case message := <-h.broadcast:

			log.Println("BROADCAST EVENT:", string(message))
			log.Println("CLIENTS COUNT:", len(h.clients))

			for client := range h.clients {
				log.Printf("sending to client id=%p", client)
				client.send <- message
			}
		}
	}
}
