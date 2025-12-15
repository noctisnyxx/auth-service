package repository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/noctisnyxx/auth-service/internal/domain"
)

type TokenRepository struct {
	logger *slog.Logger
	db     *sql.DB
}

func NewTokenRepository(logger *slog.Logger, db *sql.DB) *TokenRepository {
	return &TokenRepository{
		logger: logger,
		db:     db,
	}
}

func (r *TokenRepository) InsertRefreshToken(ctx context.Context, refreshToken domain.RefreshToken) error {
	q := `
		INSERT INTO refresh_tokens (id, token, user_id, client_id, expires_at, issued_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(
		ctx,
		q,
		refreshToken.Id,
		refreshToken.Token,
		refreshToken.UserId,
		refreshToken.ClientId,
		refreshToken.ExpiresAt,
		refreshToken.IssuedAt,
	)
	if err != nil {
		r.logger.Error("failed to insert refresh token", "error", err)
		return err
	}
	return nil

}

func (r *TokenRepository) FindRefreshToken(ctx context.Context, q domain.FindTokenQuery) ([]domain.RefreshToken, error) {
	query := `
		SELECT id, token, user_id, client_id, expires_at, issued_at, revoked_at FROM refresh_tokens WHERE 1=1
	`
	args := []interface{}{}

	if q.Token != nil {
		query += " AND token = ?"
		args = append(args, *q.Token)
	}
	if q.UserId != nil {
		query += " AND user_id = ?"
		args = append(args, *q.UserId)
	}
	if q.ClientId != nil {
		query += " AND client_id = ?"
		args = append(args, *q.ClientId)
	}
	if q.Revoked != nil {
		if *q.Revoked {
			query += " AND revoked_at IS NOT NULL"
		} else {
			query += " AND revoked_at IS NULL"
		}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("failed to find refresh tokens", "error", err)
		return nil, err
	}
	defer rows.Close()

	var tokens []domain.RefreshToken
	for rows.Next() {
		var token domain.RefreshToken
		var revokedAt sql.NullTime
		if err := rows.Scan(
			&token.Id,
			&token.Token,
			&token.UserId,
			&token.ClientId,
			&token.ExpiresAt,
			&token.IssuedAt,
			&revokedAt,
		); err != nil {
			r.logger.Error("failed to scan refresh token", "error", err)
			return nil, err
		}
		if revokedAt.Valid {
			token.RevokedAt = revokedAt.Time
		}
		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("rows error", "error", err)
		return nil, err
	}

	return tokens, nil
}

func (r *TokenRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	q := `
		UPDATE refresh_tokens SET revoked_at = CURRENT_TIMESTAMP WHERE token = ?
	`
	_, err := r.db.ExecContext(ctx, q, token)
	if err != nil {
		r.logger.Error("failed to revoke refresh token", "error", err)
		return err
	}
	return nil
}
