package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/noctisnyxx/auth-service/internal/domain"
)

type ClientRepository struct {
	logger *slog.Logger
	db     *sql.DB
}

func NewClientRepository(logger *slog.Logger, db *sql.DB) *ClientRepository {
	return &ClientRepository{
		logger: logger,
		db:     db,
	}
}

func (r *ClientRepository) Insert(ctx context.Context, client *domain.Client) error {
	r.logger.Debug("Call insert client")
	q := fmt.Sprintf("INSERT INTO %s (id, name, protocol_type_id, secret) VALUES(?,?,?,?)", TABLE_CLIENTS)
	res, err := r.db.ExecContext(
		ctx,
		q,
		client.ID,
		client.Name,
		client.ProtocolTypeID,
		client.Secret,
	)
	if err != nil {
		r.logger.Error("failed to insert client", "error", err)
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		r.logger.Error("failed to get rows affected", "error", err)
		return err
	}
	if rowsAffected == 0 {
		r.logger.Warn("no rows affected when inserting client")
	}

	return nil
}
func (r *ClientRepository) Find(ctx context.Context, q domain.FindClientQuery) ([]domain.Client, error) {
	r.logger.Debug("Call find clients")
	query := fmt.Sprintf(`
		SELECT id, name, protocol_type_id, secret, created_at, updated_at FROM %s WHERE 1=1
	`, TABLE_CLIENTS)
	args := []interface{}{}

	if q.ID != nil {
		query += " AND id = ?"
		args = append(args, *q.ID)
	}
	if q.Name != nil {
		query += " AND name = ?"
		args = append(args, *q.Name)
	}
	if q.ProtocolTypeID != nil {
		query += " AND protocol_type_id = ?"
		args = append(args, *q.ProtocolTypeID)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("failed to find clients", "error", err)
		return nil, err
	}
	defer rows.Close()
	clients := []domain.Client{}
	for rows.Next() {
		var client domain.Client
		if err := rows.Scan(
			&client.ID,
			&client.Name,
			&client.ProtocolTypeID,
			&client.Secret,
			&client.CreatedAt,
			&client.UpdatedAt,
		); err != nil {
			r.logger.Error("failed to scan client", "error", err)
			return nil, err
		}
		clients = append(clients, client)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("rows error", "error", err)
		return nil, err
	}

	return clients, nil
}
func (r *ClientRepository) Update(ctx context.Context, id string, u domain.ClientUpdate) error {
	r.logger.Debug("Call update client")
	query := fmt.Sprintf("UPDATE %s SET", TABLE_CLIENTS)
	args := []interface{}{}
	if u.Name != nil {
		query += " name = ?,"
		args = append(args, *u.Name)
	}
	if u.ProtocolTypeID != nil {
		query += " protocol_type_id = ?,"
		args = append(args, *u.ProtocolTypeID)
	}

	if len(args) == 0 {
		r.logger.Warn("no fields to update for client")
		return nil
	}
	query = query[:len(query)-1]
	query += " WHERE id = ?"
	args = append(args, id)

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("failed to update client", "error", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		r.logger.Error("failed to get rows affected", "error", err)
		return err
	}
	if rowsAffected == 0 {
		r.logger.Warn("no rows affected when updating client", "id", id)
	}

	return nil

}
func (r *ClientRepository) Delete(ctx context.Context, id string) error {
	r.logger.Debug("Call delete client")
	q := fmt.Sprintf("DELETE FROM %s WHERE id = ?", TABLE_CLIENTS)
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		r.logger.Error("failed to delete client", "error", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		r.logger.Error("failed to get rows affected", "error", err)
		return err
	}
	if rowsAffected == 0 {
		r.logger.Warn("no rows affected when deleting client", "id", id)
	}
	return nil
}
