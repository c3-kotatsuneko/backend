package repository

import (
	"context"

	"github.com/c3-kotatsuneko/backend/internal/domain/entity"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	GetUserByName(ctx context.Context, name string) (*entity.User, error)
}
