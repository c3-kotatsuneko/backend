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
	defer close(errCh)
	go func() {
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

L:
	for {
		select {
		case <-doneCh:
			break L
		default:
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				errCh <- err

			}
			fmt.Println("messageType: ", messageType)
			switch messageType {
			case websocket.CloseInvalidFramePayloadData, websocket.CloseInternalServerErr, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseMandatoryExtension, websocket.CloseMessage, websocket.CloseMessageTooBig, websocket.CloseNoStatusReceived, websocket.CloseNormalClosure, websocket.ClosePolicyViolation, websocket.CloseProtocolError, websocket.CloseServiceRestart, websocket.CloseTLSHandshake, websocket.CloseTryAgainLater, websocket.CloseUnsupportedData:
				doneCh <- struct{}{}
			case websocket.TextMessage:
				var msg any
				if err := json.Unmarshal(p, &msg); err != nil {
					errCh <- err

				}
			case websocket.BinaryMessage:
				var msg rpc.GameStatusRequest
				if err := proto.Unmarshal(p, &msg); err != nil {
					fmt.Println("err: ", err)
					errCh <- err

				}
				switch msg.Mode {
				case resources.Mode_MODE_TIME_ATTACK:
				case resources.Mode_MODE_MULTI:
					switch msg.Event {
					case resources.Event_EVENT_ENTER_ROOM:
						if err := s.eventServise.EnterRoom(ctx, msg.RoomId, msg.Player, conn); err != nil {
							errCh <- err

						}
					case resources.Event_EVENT_GAME_START:
						if err := s.eventServise.GameStart(ctx, msg.RoomId); err != nil {
							errCh <- err

						}
						go func() {
							// s.eventServise.CountDown(ctx, errCh, doneCh, msg.RoomId)
							s.eventServise.Timer(ctx, errCh, doneCh, msg.RoomId)
						}()
					case resources.Event_EVENT_STATS:
						if err := s.eventServise.Stats(ctx, msg.RoomId, msg.Player); err != nil {
							errCh <- err

						}
					case resources.Event_EVENT_STACK_BLOCK:
						if err := s.eventServise.StackBlock(ctx, msg.RoomId, msg.Player); err != nil {
							errCh <- err

						}
					case resources.Event_EVENT_RESULT:
						if err := s.eventServise.Result(ctx, msg.RoomId); err != nil {
							errCh <- err

						}
					default:
						fmt.Println("unhandling event")
						errCh <- errors.New("unhandling event")

					}
				case resources.Mode_MODE_TRAINING:
				case resources.Mode_MODE_UNKNOWN:
					errCh <- errors.New("unknown mode")

				default:
					errCh <- errors.New("unhandling mode")

				}
			default:
				fmt.Println("unhandling message type")
				errCh <- nil

			}
		}

	}
	return nil

}
