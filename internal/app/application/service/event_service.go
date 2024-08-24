package service

import (
	"context"
	"fmt"
	"time"

	"github.com/c3-kotatsuneko/backend/internal/domain/service"
	"github.com/c3-kotatsuneko/protobuf/gen/game/resources"
	"github.com/c3-kotatsuneko/protobuf/gen/game/rpc"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type IEventService interface {
	EnterRoom(ctx context.Context, roomID string, playerId *resources.Player, conn *websocket.Conn) error
	GameStart(ctx context.Context, roomID string) error
	CountDonw(ctx context.Context, roomID string) error
	Timer(ctx context.Context, roomID string) error
	Stats(ctx context.Context, roomID string, player *resources.Player) error
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
	s.msgSender.Register(roomID, player, conn, nil)
	p, err := s.msgSender.GetPlayersInRoom(roomID)
	if err != nil {
		return err
	}
	r := &rpc.GameStatusResponse{
		RoomId:  roomID,
		Event:   resources.Event_EVENT_ENTER_ROOM,
		Players: p,
		Time:    -1,
	}
	fmt.Println("response: ", r)
	data, err := proto.Marshal(r)
	if err != nil {
		return err
	}
	s.msgSender.Broadcast(ctx, roomID, data)

	return nil
}

func (s *EventService) GameStart(ctx context.Context, roomID string) error {
	p, err := s.msgSender.GetPlayersInRoom(roomID)
	if err != nil {
		return err
	}
	r := &rpc.GameStatusResponse{
		RoomId:  roomID,
		Event:   resources.Event_EVENT_GAME_START,
		Players: p,
		Time:    -1,
	}
	fmt.Println("response: ", r)
	data, err := proto.Marshal(r)
	if err != nil {
		return err
	}
	s.msgSender.Broadcast(ctx, roomID, data)
	return nil
}

func (s *EventService) CountDonw(ctx context.Context, roomID string) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timer := time.NewTimer(3 * time.Second)

	startTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		case t := <-ticker.C:
			p, err := s.msgSender.GetPlayersInRoom(roomID)
			if err != nil {
				return err
			}
			elapsedTime := t.Sub(startTime)
			r := &rpc.GameStatusResponse{
				RoomId:  roomID,
				Event:   resources.Event_EVENT_TIMER,
				Players: p,
				Time:    int32(elapsedTime.Seconds()) - 4,
			}
			fmt.Println("response: ", r)
			data, err := proto.Marshal(r)
			if err != nil {
				return err
			}
			s.msgSender.Broadcast(ctx, roomID, data)
		case <-timer.C:
			return nil
		}
	}
}

func (s *EventService) Timer(ctx context.Context, roomID string) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timer := time.NewTimer(30 * time.Second)

	startTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			return nil
		case t := <-ticker.C:
			p, err := s.msgSender.GetPlayersInRoom(roomID)
			if err != nil {
				return err
			}
			elapsedTime := t.Sub(startTime)
			r := &rpc.GameStatusResponse{
				RoomId:  roomID,
				Event:   resources.Event_EVENT_TIMER,
				Players: p,
				Time:    int32(elapsedTime.Seconds()),
			}
			fmt.Println("response: ", r)
			data, err := proto.Marshal(r)
			if err != nil {
				return err
			}
			s.msgSender.Broadcast(ctx, roomID, data)
		case <-timer.C:
			return nil
		}
	}
}

func (s *EventService) Stats(ctx context.Context, roomID string, player *resources.Player) error {
	s.msgSender.UpdatePlayer(player)
	p, err := s.msgSender.GetPlayersInRoom(roomID)
	if err != nil {
		return err
	}
	r := &rpc.GameStatusResponse{
		RoomId:  roomID,
		Event:   resources.Event_EVENT_STATS,
		Players: p,
		Time:    -1,
	}
	fmt.Println("response: ", r)
	data, err := proto.Marshal(r)
	if err != nil {
		return err
	}
	s.msgSender.Broadcast(ctx, roomID, data)
	return nil
}

func (s *EventService) Result(ctx context.Context, roomID string) error {
	panic("unimplemented")
}

func (s *EventService) ExitRoom(ctx context.Context, roomID string) error {
	panic("unimplemented")
}
