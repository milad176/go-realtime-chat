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

func (r *RoomRepository) GetAll() ([]models.Room, error) {
	query := `
		SELECT id, name, created_at
		FROM rooms
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var rooms []models.Room

	for rows.Next() {
		var room models.Room

		err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		rooms = append(rooms, room)
	}

	return rooms, nil
}
