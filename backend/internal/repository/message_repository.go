package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/milad176/go-realtime-chat/backend/internal/models"
)

type MessageRepository struct {
	db *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (r *MessageRepository) Create(roomName string, username string, content string) error {

	query := `
		INSERT INTO messages
		(id, room_name, username, content)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		uuid.New(),
		roomName,
		username,
		content,
	)

	return err
}

func (r *MessageRepository) GetByRoom(roomName string) ([]models.Message, error) {

	query := `
		SELECT
			id,
			room_name,
			username,
			content,
			created_at
		FROM messages
		WHERE room_name = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(
		context.Background(),
		query,
		roomName,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var messages []models.Message

	for rows.Next() {

		var msg models.Message

		err := rows.Scan(
			&msg.ID,
			&msg.RoomName,
			&msg.Username,
			&msg.Content,
			&msg.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}

	return messages, nil
}
