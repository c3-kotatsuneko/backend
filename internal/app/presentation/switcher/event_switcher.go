package switcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/c3-kotatsuneko/backend/internal/app/application/service"
	domainService "github.com/c3-kotatsuneko/backend/internal/domain/service"
	"github.com/c3-kotatsuneko/protobuf/gen/game/resources"
	"github.com/c3-kotatsuneko/protobuf/gen/game/rpc"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type IEventSwitcher interface {
	ISwitcher
}

type EventSwitcher struct {
	eventServise service.IEventService
	msgSender    domainService.IMessageSender
}

func NewEventSwitcher(eventServise service.IEventService, msgSender domainService.IMessageSender) IEventSwitcher {
	return &EventSwitcher{
		eventServise: eventServise,
		msgSender:    msgSender,
	}
}

func (s *EventSwitcher) Switch(ctx context.Context, conn *websocket.Conn) error {
	//  エラーを処理するチャネル
	errCh := make(chan error)
	// 強制終了したい場合は、doneCh<-struct{}{}を呼ぶ
	doneCh := make(chan struct{})
	defer close(doneCh)
	go func() {
		defer close(errCh)
		for {
			select {
			case err := <-errCh:
				if err != nil {
					fmt.Println("timer error: ", err)
					return
				}
			default:
				// do nothing
			}
		}
	}()
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		fmt.Println("messageType: ", messageType)
		switch messageType {
		case websocket.TextMessage:
			var msg any
			if err := json.Unmarshal(p, &msg); err != nil {
				return err
			}
		case websocket.BinaryMessage:
			var msg rpc.GameStatusRequest
			if err := proto.Unmarshal(p, &msg); err != nil {
				fmt.Println("err: ", err)
				return err
			}
			switch msg.Mode {
			case resources.Mode_MODE_TIME_ATTACK:
			case resources.Mode_MODE_MULTI:
				switch msg.Event {
				case resources.Event_EVENT_ENTER_ROOM:
					if err := s.eventServise.EnterRoom(ctx, msg.RoomId, msg.Players, conn); err != nil {
						return err
					}
				case resources.Event_EVENT_GAME_START:
					if err := s.eventServise.GameStart(ctx, msg.RoomId); err != nil {
						return err
					}
					s.eventServise.CountDonw(ctx, msg.RoomId)
					go s.eventServise.Timer(ctx, errCh, doneCh, msg.RoomId)
				case resources.Event_EVENT_STATS:
					if err := s.eventServise.Stats(ctx, msg.RoomId, msg.Players); err != nil {
						return err
					}
				default:
					fmt.Println("unhandling event")
					return errors.New("unhandling event")
				}
			case resources.Mode_MODE_TRAINING:
			case resources.Mode_MODE_UNKNOWN:
				return errors.New("unknown mode")
			default:
				return errors.New("unhandling mode")
			}
		default:
			fmt.Println("unhandling message type")
			return nil
		}

	}
}
