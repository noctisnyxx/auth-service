package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user with this email or username already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrEmailRequired      = errors.New("email is required")
	ErrUsernameRequired   = errors.New("username is required")
	ErrPasswordRequired   = errors.New("password is required")
)
