package switcher

import (
	"context"

	"github.com/gorilla/websocket"
)

type ISwitcher interface {
	Switch(ctx context.Context, doneCh chan struct{}, conn *websocket.Conn) error
}
