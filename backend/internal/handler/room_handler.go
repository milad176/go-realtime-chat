package handler

import (
	"encoding/json"
	"net/http"

	"github.com/milad176/go-realtime-chat/backend/internal/repository"
)

type RoomHandler struct {
	repo *repository.RoomRepository
}

func NewRoomHandler(repo *repository.RoomRepository) *RoomHandler {
	return &RoomHandler{
		repo: repo,
	}
}

type CreateRoomRequest struct {
	Name string `json:"name"`
}

func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req CreateRoomRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "room name required", http.StatusBadRequest)
		return
	}

	room, err := h.repo.CreateRoom(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(room)
}

func (h *RoomHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(rooms)
}
