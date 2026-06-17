package api

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/milad176/go-realtime-chat/backend/internal/handler"
	"github.com/milad176/go-realtime-chat/backend/internal/repository"
	"github.com/milad176/go-realtime-chat/backend/internal/ws"
)

type Server struct {
	DB  *pgxpool.Pool
	Hub *ws.Hub
}

func NewServer(db *pgxpool.Pool, hub *ws.Hub) *Server {
	return &Server{
		DB:  db,
		Hub: hub,
	}
}

func (s *Server) NewHTTPServer(port string) *http.Server {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/api/health", s.handleHealth)

	// repositories
	userRepository := repository.NewUserRepository(s.DB)
	roomRepository := repository.NewRoomRepository(s.DB)
	messageRepository := repository.NewMessageRepository(s.DB)

	// handlers
	userHandler := handler.NewUserHandler(userRepository)
	roomHandler := handler.NewRoomHandler(roomRepository)
	messageHandler := handler.NewMessageHandler(messageRepository)

	// routes
	mux.HandleFunc("POST /api/users", userHandler.CreateUser)
	mux.HandleFunc("POST /api/rooms", roomHandler.CreateRoom)
	mux.HandleFunc("GET /api/rooms", roomHandler.GetRooms)
	mux.HandleFunc("GET /api/messages", messageHandler.GetMessages)
	mux.HandleFunc("/api/ws", ws.HandleWebSocket(s.Hub))

	return &http.Server{
		Addr:    ":" + port,
		Handler: cors(mux),
	}

}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {

	response := map[string]string{
		"status":   "ok",
		"database": "connected",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
