package service

import (
	"context"

	"github.com/c3-kotatsuneko/protobuf/gen/game/resources"
	"github.com/gorilla/websocket"
)

type IMessageSender interface {
	Send(ctx context.Context, to string, data interface{}) error
	Broadcast(ctx context.Context, roomID string, data interface{}) error
	Register(roomID string, player *resources.Player, conn *websocket.Conn, err chan error)
	Unregister(userID, RoomId string)
	GetPlayersInRoom(roomID string) ([]*resources.Player, error)
	UpdatePlayer(player *resources.Player) error
}
