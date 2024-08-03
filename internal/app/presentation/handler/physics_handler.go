package handler

import (
	"net/http"

	"github.com/c3-kotatsuneko/backend/internal/app/application/service"
	"github.com/c3-kotatsuneko/backend/internal/app/presentation/switcher"
	"github.com/c3-kotatsuneko/backend/internal/app/presentation/websocket"
)

type PhysicsHandler struct {
	physicsService  service.IRoomObjectService
	wsHandler       websocket.IWSHandler
	physicsSwitcher switcher.IPhysicsSwitcher
}

func NewPhysicsHandler(physicsService service.IRoomObjectService, wsHandler websocket.IWSHandler, physicsSwitcher switcher.IPhysicsSwitcher) *PhysicsHandler {
	return &PhysicsHandler{
		physicsService:  physicsService,
		wsHandler:       wsHandler,
		physicsSwitcher: physicsSwitcher,
	}
}

func (h *PhysicsHandler) Calculate() func(http.ResponseWriter, *http.Request) error {

	return func(w http.ResponseWriter, r *http.Request) error {
		h.wsHandler.Start(r.Context(), w, r, h.physicsSwitcher)
		return nil
	}
}
