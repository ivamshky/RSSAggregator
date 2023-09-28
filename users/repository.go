package users

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExist     = errors.New("row does not exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type UserRepository interface {
	Create(ctx context.Context, website User) (*User, error)
	GetById(ctx context.Context, ID uuid.UUID) (*User, error)
}
