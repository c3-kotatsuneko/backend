package handler

import (
	"net/http"

	"github.com/c3-kotatsuneko/backend/internal/app/presentation/switcher"
	"github.com/c3-kotatsuneko/backend/internal/app/presentation/websocket"
)

type PhysicsHandler struct {
	wsHandler       websocket.IWSHandler
	physicsSwitcher switcher.IPhysicsSwitcher
}

func NewPhysicsHandler(wsHandler websocket.IWSHandler, physicsSwitcher switcher.IPhysicsSwitcher) *PhysicsHandler {
	return &PhysicsHandler{
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
