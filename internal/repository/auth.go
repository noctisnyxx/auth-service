package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/noctisnyxx/auth-service/internal/domain"
)

type AuthRepository struct {
	logger *slog.Logger
	db     *sql.DB
}

func NewAuthRepository(logger *slog.Logger, db *sql.DB) *AuthRepository {
	return &AuthRepository{
		logger: logger,
		db:     db,
	}
}

func (r *AuthRepository) InsertCode(ctx context.Context, authCode domain.AuthCode) error {
	q := `
		INSERT INTO codes (code, user_id, client_id, expires_at)
		VALUES (?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(
		ctx,
		q,
		authCode.Code,
		authCode.UserID,
		authCode.ClientID,
		authCode.ExpiresAt,
	)
	if err != nil {
		r.logger.Error("failed to insert auth code", "error", err)
		return err
	}
	return nil
}

func (r *AuthRepository) FindAuthCode(ctx context.Context, code string) ([]domain.AuthCode, error) {
	q := `
		SELECT code, user_id, client_id, expires_at FROM codes WHERE code = ?
	`
	rows, err := r.db.QueryContext(ctx, q, code)
	if err != nil {
		r.logger.Error("failed to find auth code", "error", err)
		return nil, err
	}
	defer rows.Close()

	var authCodes []domain.AuthCode
	for rows.Next() {
		var authCode domain.AuthCode
		if err := rows.Scan(
			&authCode.Code,
			&authCode.UserID,
			&authCode.ClientID,
			&authCode.ExpiresAt,
		); err != nil {
			r.logger.Error("failed to scan auth code", "error", err)
			return nil, err
		}
		authCodes = append(authCodes, authCode)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("rows error", "error", err)
		return nil, err
	}

	return authCodes, nil
}

func (r *AuthRepository) DeleteAuthCode(ctx context.Context, code string) error {
	q := `
		DELETE FROM codes WHERE code = ?
	`
	_, err := r.db.ExecContext(ctx, q, code)
	if err != nil {
		r.logger.Error("failed to insert auth code", "error", err)
		return err
	}
	return nil

}
