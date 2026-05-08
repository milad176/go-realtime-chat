package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/milad176/go-realtime-chat/backend/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(username string) (*models.User, error) {
	user := &models.User{
		ID:       uuid.New().String(),
		Username: username,
	}

	query := `
		INSERT INTO users (id, username)
		VALUES ($1, $2)
		RETURNING created_at
	`

	err := r.db.QueryRow(
		context.Background(),
		query,
		user.ID,
		user.Username,
	).Scan(&user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}
