package websocket

import (
	"context"
	"net/http"

	appService "github.com/c3-kotatsuneko/backend/internal/app/application/service"
	"github.com/c3-kotatsuneko/backend/internal/app/presentation/switcher"
	"github.com/c3-kotatsuneko/backend/internal/domain/service"
	"github.com/gorilla/websocket"
)

type IWSHandler interface {
	Start(ctx context.Context, w http.ResponseWriter, r *http.Request, switcher switcher.ISwitcher) error
}

type WSHandler struct {
	roomObjectService appService.IRoomObjectService
	msgSender         service.IMessageSender
}

func NewWSHandler(roomObjectService appService.IRoomObjectService, sender service.IMessageSender) IWSHandler {
	return &WSHandler{
		roomObjectService: roomObjectService,
		msgSender:         sender,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (ws *WSHandler) Start(ctx context.Context, w http.ResponseWriter, r *http.Request, switcher switcher.ISwitcher) error {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	doneCh := make(chan struct{})

	go func() {
		if err := switcher.Switch(ctx, doneCh, conn); err != nil {
			return
		}
	}()
	<-doneCh
	// defer close(errCh)
	return nil
}
