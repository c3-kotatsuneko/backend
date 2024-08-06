package service

import (
	"context"
	"fmt"

	"github.com/c3-kotatsuneko/backend/internal/domain/service"
	"github.com/c3-kotatsuneko/protobuf/gen/game/resources"
	"github.com/c3-kotatsuneko/protobuf/gen/game/rpc"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type IEventService interface {
	EnterRoom(ctx context.Context, roomID string, playerId *resources.Player, conn *websocket.Conn) error
	GameStart(ctx context.Context, roomID string) error
	Timer(ctx context.Context, roomID string) error
	Stats(ctx context.Context, roomID string) error
	Result(ctx context.Context, roomID string) error
	ExitRoom(ctx context.Context, roomID string) error
}

type EventService struct {
	msgSender service.IMessageSender
}

func NewEventService(msgSender service.IMessageSender) IEventService {
	return &EventService{
		msgSender: msgSender,
	}
}

func (s *EventService) EnterRoom(ctx context.Context, roomID string, player *resources.Player, conn *websocket.Conn) error {
	fmt.Println("EnterRoom")
	s.msgSender.Register(roomID, player, conn, nil)
	p, err := s.msgSender.GetPlayersInRoom(roomID)
	fmt.Println("players: ", p)
	if err != nil {
		return err
	}
	fmt.Println("players: ", p)
	r := &rpc.GameStatusResponse{
		RoomId:  roomID,
		Event:   resources.Event_EVENT_ENTER_ROOM,
		Players: p,
		Time:    -1,
	}
	data, err := proto.Marshal(r)
	if err != nil {
		return err
	}
	s.msgSender.Broadcast(ctx, roomID, data)

	return nil
}

func (s *EventService) GameStart(ctx context.Context, roomID string) error {
	return nil
}

func (s *EventService) Timer(ctx context.Context, roomID string) error {
	return nil
}

func (s *EventService) Stats(ctx context.Context, roomID string) error {
	panic("unimplemented")
}

func (s *EventService) Result(ctx context.Context, roomID string) error {
	panic("unimplemented")
}

func (s *EventService) ExitRoom(ctx context.Context, roomID string) error {
	panic("unimplemented")
}
