package domain

import (
	"context"
	"time"
)

type IUserRepository interface {
	Insert(ctx context.Context, u *User) error
	Find(ctx context.Context, q FindUserQuery) ([]User, error)
	Update(ctx context.Context, id string, u UserUpdate) error
	Delete(ctx context.Context, id string) error
}

type User struct {
	ID            string
	Username      string
	Email         string
	Password      string
	Phone         *string
	IsActive      bool
	EmailVerified bool
	PhoneVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type UserUpdate struct {
	Username      *string
	Email         *string
	Password      *string
	Phone         *string
	IsActive      *bool
	EmailVerified *bool
	PhoneVerified *bool
}

type FindUserQuery struct {
	ID            *string
	Email         *string
	Username      *string
	EmailVerified *bool
	PhoneVerified *bool
}
