package presentation

import (
	"github.com/c3-kotatsuneko/backend/internal/app/presentation/handler"
)

type Root struct {
	PhysicsHandler *handler.PhysicsHandler
}

func New(physicsHandler *handler.PhysicsHandler) *Root {
	return &Root{
		PhysicsHandler: physicsHandler,
	}
}
