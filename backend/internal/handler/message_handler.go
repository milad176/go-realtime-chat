package handler

import (
	"encoding/json"
	"net/http"

	"github.com/milad176/go-realtime-chat/backend/internal/repository"
)

type MessageHandler struct {
	repo *repository.MessageRepository
}

func NewMessageHandler(repo *repository.MessageRepository) *MessageHandler {
	return &MessageHandler{
		repo: repo,
	}
}

func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {

	room := r.URL.Query().Get("room")

	messages, err := h.repo.GetByRoom(room)

	if err != nil {
		http.Error(w, "failed to fetch messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(messages)
}
