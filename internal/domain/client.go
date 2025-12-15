package domain

import (
	"context"
	"time"
)

type IClientRepository interface {
	Insert(ctx context.Context, c *Client) error
	Find(ctx context.Context, q FindClientQuery) ([]Client, error)
	Update(ctx context.Context, id string, u ClientUpdate) error
	Delete(ctx context.Context, id string) error
}

type Client struct {
	ID             string
	Name           string
	ProtocolTypeID string
	Secret         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type ProtocolType struct {
	ID   string
	Name string
}

type ClientUpdate struct {
	Name           *string
	ProtocolTypeID *string
}

type FindClientQuery struct {
	ID             *string
	Secret         *string
	Name           *string
	ProtocolTypeID *string
}
