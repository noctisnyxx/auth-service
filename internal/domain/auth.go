package domain

import (
	"context"
	"time"
)

type IAuthRepository interface {
	InsertCode(ctx context.Context, authCode AuthCode) error
	FindAuthCode(ctx context.Context, code string) ([]AuthCode, error)
	DeleteAuthCode(ctx context.Context, code string) error
}

type AuthCode struct {
	Code      string
	UserID    string
	ClientID  string
	ExpiresAt time.Time
}
