package api

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	DB *pgxpool.Pool
}

func NewServer(db *pgxpool.Pool) *Server {
	return &Server{DB: db}
}

func (s *Server) Start(port string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", s.handleHealth)

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
