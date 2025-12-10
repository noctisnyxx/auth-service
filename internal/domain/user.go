package domain

import (
	"context"
	"time"
)

type IUserRepository interface {
	Insert(context.Context, *User) error
	FindUsers(context.Context, FindUserQuery) ([]User, error)
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

type FindUserQuery struct {
	ID            *string
	Email         *string
	Username      *string
	EmailVerified *bool
	PhoneVerified *bool
}
