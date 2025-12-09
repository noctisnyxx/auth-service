package repository

import (
	"database/sql"
	"log/slog"

	"repo.polytron.co.id/rd-special-project/auth-service/backend/internal/domain"
)

/* type UserRepository interface {
	Create(user *User) error
	FindByUsername(username string) (*User, error)
	FindByEmail(email string) (*User, error)
} */

type UserRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewUserRepository(db *sql.DB, logger *slog.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepository) Create(user *domain.User) error {

	return nil
}
