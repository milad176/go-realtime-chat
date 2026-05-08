package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/milad176/go-realtime-chat/backend/internal/models"
)

type RoomRepository struct {
	db *pgxpool.Pool
}

func NewRoomRepository(db *pgxpool.Pool) *RoomRepository {

	return &RoomRepository{db: db}
}

func (r *RoomRepository) CreateRoom(name string) (*models.Room, error) {
	room := &models.Room{

		ID:   uuid.NewString(),
		Name: name,
	}

	query := `
		INSERT INTO rooms (id, name)
		VALUES ($1, $2)
		RETURNING created_at
	`
	err := r.db.QueryRow(
		context.Background(),
		query,
		room.ID,
		room.Name,
	).Scan(&room.CreatedAt)

	if err != nil {
		return nil, err
	}

	return room, nil
}
