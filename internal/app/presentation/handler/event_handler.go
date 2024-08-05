package handler

import (
	"net/http"

	"github.com/c3-kotatsuneko/backend/internal/app/presentation/switcher"
	"github.com/c3-kotatsuneko/backend/internal/app/presentation/websocket"
)

type EventHandler struct {
	wsHandler     websocket.IWSHandler
	eventSwitcher switcher.IEventSwitcher
}

func NewEventHandler(wsHandler websocket.IWSHandler, eventSwitcher switcher.IEventSwitcher) *EventHandler {
	return &EventHandler{
		wsHandler:     wsHandler,
		eventSwitcher: eventSwitcher,
	}
}

func (h *EventHandler) ManageEvent() func(http.ResponseWriter, *http.Request) error {

	return func(w http.ResponseWriter, r *http.Request) error {
		h.wsHandler.Start(r.Context(), w, r, h.eventSwitcher)
		return nil
	}
}
