package switcher

import (
	"context"
	"encoding/json"
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
			switch msg.Event {
			case resources.Event_EVENT_TIME_ATTACK:
				if err := s.eventServise.TimeAttack(ctx, msg.RoomId, msg.RoomId); err != nil {
					return err
				}
			case resources.Event_EVENT_MULTI:
				if err := s.eventServise.MultiPlay(ctx, msg.RoomId, msg.RoomId); err != nil {
					return err
				}
			case resources.Event_EVENT_TRAINING:
				if err := s.eventServise.Training(ctx, msg.RoomId, msg.RoomId); err != nil {
					return err
				}
			default:
				return nil
			}
		default:
			return nil
		}

	}
}
