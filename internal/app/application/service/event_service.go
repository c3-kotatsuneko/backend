package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/c3-kotatsuneko/backend/internal/app/constants"
	"github.com/c3-kotatsuneko/backend/internal/domain/service"
	"github.com/c3-kotatsuneko/protobuf/gen/game/resources"
	"github.com/c3-kotatsuneko/protobuf/gen/game/rpc"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type IEventService interface {
	EnterRoom(ctx context.Context, roomID string, playerId *resources.Player, conn *websocket.Conn) error
	GameStart(ctx context.Context, roomID string) error
	CountDown(ctx context.Context, timerCh chan<- error, doneCh <-chan struct{}, roomID string)
	Timer(ctx context.Context, timerCh chan<- error, doneCh <-chan struct{}, roomID string)
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
	isRegistered := s.msgSender.IsPlayerRegistered(player.PlayerId)
	if isRegistered {
		fmt.Println("player is already registered")
		return errors.New("player is already registered")
	}
	status, _ := s.msgSender.GetRoomStatus(roomID)
	if status == constants.RoomStatusPlaying {
		fmt.Println("room is playing1")
		return errors.New("room is playing1")
	}
	p, err := s.msgSender.GetPlayersInRoom(roomID)
	if err != nil {
		return err
	}
	i := len(p)
	if i < len(constants.Directions) {
		player.Direction = constants.Directions[i]
	} else {
		return errors.New("room is full")
	}
	s.msgSender.Register(roomID, player, conn, nil)
	s.msgSender.SetRoomStatus(roomID, constants.RoomStatusWaiting)
	players, err := s.msgSender.GetPlayersInRoom(roomID)
	if err != nil {
		return err
	}
	r := &rpc.GameStatusResponse{
		RoomId:  roomID,
		Event:   resources.Event_EVENT_ENTER_ROOM,
		Players: players,
		Time:    -1,
		Mode:    resources.Mode_MODE_MULTI,
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
	status, err := s.msgSender.GetRoomStatus(roomID)
	fmt.Println("status: ", status)
	if err != nil {
		return err
	}
	if status == constants.RoomStatusPlaying {
		fmt.Println("room is playing2")
		return errors.New("room is playing2")
	}
	s.msgSender.SetRoomStatus(roomID, constants.RoomStatusPlaying)
	p, err := s.msgSender.GetPlayersInRoom(roomID)
	if err != nil {
		return err
	}
	r := &rpc.GameStatusResponse{
		RoomId:  roomID,
		Event:   resources.Event_EVENT_GAME_START,
		Players: p,
		Time:    -1,
		Mode:    resources.Mode_MODE_MULTI,
	}
	fmt.Println("response: ", r)
	data, err := proto.Marshal(r)
	if err != nil {
		return err
	}
	s.msgSender.Broadcast(ctx, roomID, data)
	return nil
}

func (s *EventService) CountDown(ctx context.Context, timerCh chan<- error, doneCh <-chan struct{}, roomID string) {
	ticker := time.NewTicker(time.Duration(constants.IntervalTicker) * time.Second)
	defer ticker.Stop()

	timer := time.NewTimer(time.Duration(constants.CountDownTimer) * time.Second)

	startTime := -1*constants.CountDownTimer - 1

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p, err := s.msgSender.GetPlayersInRoom(roomID)
			if err != nil {
				timerCh <- err
			}
			startTime++
			r := &rpc.GameStatusResponse{
				RoomId:  roomID,
				Event:   resources.Event_EVENT_TIMER,
				Players: p,
				Time:    int32(startTime),
				Mode:    resources.Mode_MODE_MULTI,
			}
			fmt.Println("response: ", r)
			data, err := proto.Marshal(r)
			if err != nil {
				timerCh <- err
			}
			s.msgSender.Broadcast(ctx, roomID, data)
		case <-timer.C:
			return
		case <-doneCh:
			return
		}
	}
}

func (s *EventService) Timer(ctx context.Context, timerCh chan<- error, doneCh <-chan struct{}, roomID string) {
	ticker := time.NewTicker(time.Duration(constants.IntervalTicker) * time.Second)
	defer ticker.Stop()

	timer := time.NewTimer(time.Duration(constants.TimeOutTimer) * time.Second)

	startTime := constants.TimeOutTimer + 1

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p, err := s.msgSender.GetPlayersInRoom(roomID)
			if err != nil {
				timerCh <- err
			}
			startTime--
			r := &rpc.GameStatusResponse{
				RoomId:  roomID,
				Event:   resources.Event_EVENT_TIMER,
				Players: p,
				Time:    int32(startTime),
				Mode:    resources.Mode_MODE_MULTI,
			}
			fmt.Println("response: ", r)
			data, err := proto.Marshal(r)
			if err != nil {
				timerCh <- err
			}
			s.msgSender.Broadcast(ctx, roomID, data)
		case <-timer.C:
			return
		case <-doneCh:
			return
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
		Mode:    resources.Mode_MODE_MULTI,
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
