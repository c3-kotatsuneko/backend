package switcher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/c3-kotatsuneko/backend/internal/app/application/service"
	"github.com/c3-kotatsuneko/backend/internal/domain/entity"
	domainService "github.com/c3-kotatsuneko/backend/internal/domain/service"
	"github.com/c3-kotatsuneko/protobuf/gen/game/resources"
	"github.com/c3-kotatsuneko/protobuf/gen/game/rpc"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type IPhysicsSwitcher interface {
	ISwitcher
}

type PhysicsSwitcher struct {
	physicsService service.IRoomObjectService
	msgSender      domainService.IMessageSender
}

func NewPhysicsSwitcher(physicsService service.IRoomObjectService, msgSender domainService.IMessageSender) IPhysicsSwitcher {
	return &PhysicsSwitcher{
		physicsService: physicsService,
		msgSender:      msgSender,
	}
}

func (s *PhysicsSwitcher) Switch(ctx context.Context, doneCh chan struct{}, conn *websocket.Conn) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Recovered from a panic: %v", err)
		}
	}()
	// x := &rpc.PhysicsRequest{
	// 	SenderId: "roomID",
	// 	RoomId:   "roomID",
	// 	Hands: &resources.Hand{
	// 		UserId: "roomID",
	// 		State:  1,
	// 		CenterPosition: &resources.Vector3{
	// 			X: 0,
	// 			Y: 0,
	// 			Z: 0,
	// 		},
	// 		ActionPosition: &resources.Vector3{
	// 			X: 0,
	// 			Y: 0,
	// 			Z: 0,
	// 		},
	// 	},
	// }
	// b, err := proto.Marshal(x)
	// if err != nil {
	// 	fmt.Println("err: ", err)
	// 	return err
	// }
	// if err := s.msgSender.Send(ctx, "roomID", b); err != nil {
	// 	fmt.Println("err: ", err)
	// 	return err
	// }

	errCh := make(chan error)
	messageType, p, err := conn.ReadMessage()
	if err != nil {
		errCh <- err
	}
	fmt.Println("messageType: ", messageType)
	switch messageType {
	case websocket.TextMessage:
		var msg any
		if err := json.Unmarshal(p, &msg); err != nil {
			errCh <- err
		}
	case websocket.BinaryMessage:
		var msg rpc.PhysicsRequest
		if err := proto.Unmarshal(p, &msg); err != nil {
			fmt.Println("err: ", err)
			errCh <- err
		}
		// fmt.Println("msg: ", msg)
		// hand := &entity.Hand{
		// 	UserID: msg.SenderId,
		// 	State:  entity.HandState(msg.Hands.State),
		// 	CenterPosition: entity.Vector3{
		// 		X: msg.Hands.CenterPosition.X,
		// 		Y: msg.Hands.CenterPosition.Y,
		// 		Z: msg.Hands.CenterPosition.Z,
		// 	},
		// 	ActionPosition: entity.Vector3{
		// 		X: msg.Hands.ActionPosition.X,
		// 		Y: msg.Hands.ActionPosition.Y,
		// 		Z: msg.Hands.ActionPosition.Z,
		// 	},
		// }
		// fmt.Println("000000000000")
		ctx := context.WithValue(ctx, "roomID", msg.RoomId)
		player := &resources.Player{
			PlayerId: msg.SenderId,
		}

		s.msgSender.Register(msg.RoomId, player, conn, errCh)
		s.physicsService.Init(ctx, msg.RoomId)
	default:
	}
	go s.readPump(ctx, conn, errCh)
	go s.writePump(ctx, errCh)
	<-ctx.Done()
	return nil
}

func (s *PhysicsSwitcher) readPump(ctx context.Context, conn *websocket.Conn, errCh chan<- error) ([]byte, error) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			errCh <- err
		}
		fmt.Println("messageType: ", messageType)
		switch messageType {
		case websocket.TextMessage:
			var msg any
			if err := json.Unmarshal(p, &msg); err != nil {
				errCh <- err
			}
		case websocket.BinaryMessage:
			var msg rpc.PhysicsRequest
			if err := proto.Unmarshal(p, &msg); err != nil {
				fmt.Println("err: ", err)
				errCh <- err
			}
			// fmt.Println("msg: ", msg)
			// hand := &entity.Hand{
			// 	UserID: msg.SenderId,
			// 	State:  entity.HandState(msg.Hands.State),
			// 	CenterPosition: entity.Vector3{
			// 		X: msg.Hands.CenterPosition.X,
			// 		Y: msg.Hands.CenterPosition.Y,
			// 		Z: msg.Hands.CenterPosition.Z,
			// 	},
			// 	ActionPosition: entity.Vector3{
			// 		X: msg.Hands.ActionPosition.X,
			// 		Y: msg.Hands.ActionPosition.Y,
			// 		Z: msg.Hands.ActionPosition.Z,
			// 	},
			// }
			// fmt.Println("000000000000")
			objs := make([]*entity.Object, 0, len(msg.Objects))
			for _, obj := range msg.Objects {
				objs = append(objs, &entity.Object{
					ID:    obj.ObjectId,
					Layer: obj.Layer,
					Kinds: entity.ObjectKind(obj.Kinds),
					State: entity.ObjectState(obj.State),
					Position: entity.Vector3{
						X: obj.Position.X,
						Y: obj.Position.Y,
						Z: obj.Position.Z,
					},
					Size: entity.Vector3{
						X: obj.Size.X,
						Y: obj.Size.Y,
						Z: obj.Size.Z,
					},
				})
			}
			if err := s.physicsService.Share(ctx, msg.RoomId, objs); err != nil {
				errCh <- err
			}
		default:
		}

	}
}

func (s *PhysicsSwitcher) writePump(ctx context.Context, errCh chan error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	roomID := ctx.Value("roomID").(string)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-errCh:
			fmt.Println("err: ", err)
			cancel()
			//TODO: エラー処理
			return
		case <-ticker.C:
			objs, err := s.physicsService.Get(ctx, roomID)
			if err != nil {
				fmt.Println("err: ", err)
				errCh <- err
				return
			}
			resourceObj := make([]*resources.Object, 0, len(objs))
			for _, obj := range objs {
				resourceObj = append(resourceObj, &resources.Object{
					ObjectId: obj.ID,
					Layer:    obj.Layer,
					Kinds:    resources.ObjectKind(obj.Kinds),
					State:    resources.ObjectState(obj.State),
					Position: &resources.Vector3{
						X: obj.Position.X,
						Y: obj.Position.Y,
						Z: obj.Position.Z,
					},
					Size: &resources.Vector3{
						X: obj.Size.X,
						Y: obj.Size.Y,
						Z: obj.Size.Z,
					},
				})
			}
			fmt.Println("resourceObj", resourceObj)
			physics := rpc.PhysicsResponse{
				RoomId:   roomID,
				SenderId: "SenderId",
				Objects:  resourceObj,
			}

			b, err := proto.Marshal(&physics)
			if err != nil {
				fmt.Println(err)
				errCh <- err
			}
			fmt.Println("send physics")
			if err := s.msgSender.Broadcast(ctx, roomID, b); err != nil {
				errCh <- err
			}
		}
	}
}
