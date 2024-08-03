package repository

import (
	"context"

	"github.com/c3-kotatsuneko/backend/internal/domain/entity"
)

type IRoomObjectRepository interface {
	Get(ctx context.Context, roomID string) (*[]entity.Object, error)
	Set(ctx context.Context, roomID string, objects *[]entity.Object) error
	Resister(ctx context.Context, id string) error
	Unregister(ctx context.Context, id string) error
}
