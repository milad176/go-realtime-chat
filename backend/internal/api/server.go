package api

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/milad176/go-realtime-chat/backend/internal/handler"
	"github.com/milad176/go-realtime-chat/backend/internal/repository"
)

type Server struct {
	DB *pgxpool.Pool
}

func NewServer(db *pgxpool.Pool) *Server {
	return &Server{DB: db}
}

func (s *Server) Start(port string) error {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/api/health", s.handleHealth)

	// repositories
	userRepository := repository.NewUserRepository(s.DB)
	roomRepository := repository.NewRoomRepository(s.DB)

	// handlers
	userHandler := handler.NewUserHandler(userRepository)
	roomHandler := handler.NewRoomHandler(roomRepository)

	// routes
	mux.HandleFunc("POST /api/users", userHandler.CreateUser)
	mux.HandleFunc("POST /api/rooms", roomHandler.CreateRoom)
	mux.HandleFunc("GET /api/rooms", roomHandler.GetRooms)

	return http.ListenAndServe(":"+port, mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {

	response := map[string]string{
		"status":   "ok",
		"database": "connected",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
