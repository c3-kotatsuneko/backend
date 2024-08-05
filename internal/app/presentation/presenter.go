package presentation

import (
	"github.com/c3-kotatsuneko/backend/internal/app/presentation/handler"
)

type Root struct {
	PhysicsHandler *handler.PhysicsHandler
	EventHandler   *handler.EventHandler
}

func New(physicsHandler *handler.PhysicsHandler, EventHandler *handler.EventHandler) *Root {
	return &Root{
		PhysicsHandler: physicsHandler,
		EventHandler:   EventHandler,
	}
}
