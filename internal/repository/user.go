package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/noctisnyxx/auth-service/internal/domain"
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

func (r *UserRepository) Insert(ctx context.Context, user *domain.User) error {
	r.logger.Debug("Call insert user")
	q := fmt.Sprintf(`
		INSERT INTO %s (id, username, email, phone, password) VALUES(?, ?, ?, ?, ?)
	`, TABLE_USERS)
	res, err := r.db.ExecContext(
		ctx,
		q,
		user.ID,
		user.Username,
		user.Email,
		user.Phone,
		user.Password,
	)
	if err != nil {
		r.logger.Error("failed to insert user", "error", err)
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		r.logger.Error("failed to get rows affected", "error", err)
		return err
	}
	if rowsAffected == 0 {
		r.logger.Warn("no rows affected when inserting user")
	}
	return nil
}

func (r *UserRepository) FindUsers(ctx context.Context, q domain.FindUserQuery) ([]domain.User, error) {
	r.logger.Debug("Call find users")
	query := fmt.Sprintf(`
		SELECT id, username, email, phone, password, is_active, email_verified, phone_verified, created_at, updated_at FROM %s WHERE 1=1
	`, TABLE_USERS)
	args := []interface{}{}

	if q.ID != nil {
		query += " AND id = ?"
		args = append(args, *q.ID)
	}
	if q.Email != nil {
		query += " AND email = ?"
		args = append(args, *q.Email)
	}
	if q.Username != nil {
		query += " AND username = ?"
		args = append(args, *q.Username)
	}
	if q.EmailVerified != nil {
		query += " AND email_verified = ?"
		args = append(args, *q.EmailVerified)
	}
	if q.PhoneVerified != nil {
		query += " AND phone_verified = ?"
		args = append(args, *q.PhoneVerified)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("failed to find users", "error", err)
		return nil, err
	}
	defer rows.Close()
	users := []domain.User{}
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Phone,
			&user.Password,
			&user.IsActive,
			&user.EmailVerified,
			&user.PhoneVerified,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			r.logger.Error("failed to scan user", "error", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("rows error", "error", err)
		return nil, err
	}

	return users, nil

}
