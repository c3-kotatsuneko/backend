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

func (s *EventSwitcher) Switch(ctx context.Context, doneCh chan struct{}, conn *websocket.Conn) error {
	//  エラーを処理するチャネル
	errCh := make(chan error)
	// 終了したい場合は、doneCh<-struct{}{}を呼ぶ
	// doneCh := make(chan struct{})
	// defer close(doneCh)
	go func() {
		defer close(errCh)
		for {
			select {
			case err := <-errCh:
				if err != nil {
					fmt.Println("timer error: ", err)
					return
				}
				doneCh <- struct{}{}
			default:
				// do nothing
			}
		}
	}()
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			errCh <- err
			break

		}
		fmt.Println("messageType: ", messageType)
		switch messageType {
		case websocket.TextMessage:
			var msg any
			if err := json.Unmarshal(p, &msg); err != nil {
				errCh <- err
				break

			}
		case websocket.BinaryMessage:
			var msg rpc.GameStatusRequest
			if err := proto.Unmarshal(p, &msg); err != nil {
				fmt.Println("err: ", err)
				errCh <- err
				break

			}
			switch msg.Mode {
			case resources.Mode_MODE_TIME_ATTACK:
			case resources.Mode_MODE_MULTI:
				switch msg.Event {
				case resources.Event_EVENT_ENTER_ROOM:
					if err := s.eventServise.EnterRoom(ctx, msg.RoomId, msg.Player, conn); err != nil {
						errCh <- err
						break

					}
				case resources.Event_EVENT_GAME_START:
					if err := s.eventServise.GameStart(ctx, msg.RoomId); err != nil {
						errCh <- err
						break

					}
					go func() {
						// s.eventServise.CountDown(ctx, errCh, doneCh, msg.RoomId)
						s.eventServise.Timer(ctx, errCh, doneCh, msg.RoomId)
					}()
				case resources.Event_EVENT_STATS:
					if err := s.eventServise.Stats(ctx, msg.RoomId, msg.Player); err != nil {
						errCh <- err
						break

					}
				case resources.Event_EVENT_STACK_BLOCK:
					if err := s.eventServise.StackBlock(ctx, msg.RoomId, msg.Player); err != nil {
						errCh <- err
						break

					}
				case resources.Event_EVENT_RESULT:
					if err := s.eventServise.Result(ctx, msg.RoomId); err != nil {
						errCh <- err
						break

					}
				default:
					fmt.Println("unhandling event")
					errCh <- errors.New("unhandling event")
					break

				}
			case resources.Mode_MODE_TRAINING:
			case resources.Mode_MODE_UNKNOWN:
				errCh <- errors.New("unknown mode")
				break

			default:
				errCh <- errors.New("unhandling mode")
				break

			}
		default:
			fmt.Println("unhandling message type")
			errCh <- nil
			break
		}
	}
	doneCh <- struct{}{}
	return nil
}
