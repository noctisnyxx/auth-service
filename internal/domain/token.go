package domain

import (
	"context"
	"time"
)

type ITokenRepository interface {
	InsertRefreshToken(ctx context.Context, refreshToken RefreshToken) error
	FindRefreshToken(ctx context.Context, q FindTokenQuery) ([]RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
}
type TokenClaims struct {
	UserId         string
	ClientId       string
	ExpireDuration time.Duration
}

type Token struct {
	AccessToken  string
	RefreshToken string
}

type RefreshToken struct {
	Id        string
	Token     string
	UserId    string
	ClientId  string
	ExpiresAt time.Time
	IssuedAt  time.Time
	RevokedAt time.Time
}

type FindTokenQuery struct {
	Token    *string
	UserId   *string
	ClientId *string
	Revoked  *bool
}
