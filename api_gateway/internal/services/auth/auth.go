package auth

import (
	"context"
	"errors"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInternalServer     = errors.New("internal server error")
)

type Service interface {
	Register(ctx context.Context, login, password string) error
	Login(ctx context.Context, login, password string, appID uint32) (token string, err error)
	IsAdmin(ctx context.Context, userID uint64, appID uint32) (isAdmin bool, err error)
}
