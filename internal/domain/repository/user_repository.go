package repository

import (
	"github.com/c3-kotatsuneko/backend/internal/domain/entity"
)

type IUserRepository interface {
	CreateUser(user *entity.User) error
	GetUserByID(id string) (*entity.User, error)
	GetUserByName(name string) (*entity.User, error)
}
