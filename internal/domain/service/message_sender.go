package service

import (
	"context"

	"github.com/c3-kotatsuneko/protobuf/gen/game/resources"
	"github.com/gorilla/websocket"
)

type IMessageSender interface {
	Send(ctx context.Context, to string, data interface{}) error
	Broadcast(ctx context.Context, roomID string, data interface{}) error
	IsPlayerRegistered(playerID string) bool
	Register(roomID string, player *resources.Player, conn *websocket.Conn, err chan error)
	Unregister(userID, RoomId string)
	GetPlayersInRoom(roomID string) ([]*resources.Player, error)
	UpdatePlayer(player *resources.Player) error
	SetRoomStatus(roomID string, status string) error
	GetRoomStatus(roomID string) (status string, err error)
	GetTime(ctx context.Context, roomID string) int32
	StartTimer(ctx context.Context, doneCh <-chan struct{}, roomID string)
	DestroyRoom(ctx context.Context, roomID string)
}
