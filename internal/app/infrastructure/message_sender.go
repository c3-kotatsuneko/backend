package infrastructure

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/c3-kotatsuneko/backend/internal/app/constants"
	"github.com/c3-kotatsuneko/backend/internal/domain/service"
	"github.com/c3-kotatsuneko/protobuf/gen/game/resources"
	"github.com/c3-kotatsuneko/protobuf/gen/game/rpc"
	"github.com/gorilla/websocket"
)

const writeWait = 500 * time.Millisecond

type Client struct {
	conn   *websocket.Conn
	cancel chan struct{}
	ch     chan interface{}
	err    chan error
}

func (c *Client) run() {
	for {
		select {
		case <-c.cancel:
			return
		case msg := <-c.ch:
			switch msg := msg.(type) {
			case []byte:
				fmt.Println("1111")
				err := c.conn.WriteMessage(websocket.BinaryMessage, msg)
				if err != nil {
					c.err <- err
					return
				}
			case string:
				fmt.Println("22222")

				err := c.conn.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					c.err <- err
					return
				}
			default:
				fmt.Println(msg)

				c.err <- errors.New("unknown message type")
			}
		}
	}
}

type MsgSender struct {
	mutex   *sync.RWMutex
	clients map[string]*Client           // userID -> client
	players map[string]*resources.Player // playerID -> Player
	rooms   map[string]*Room             // roomID -> {userIDs, status}
}

type Room struct {
	userIDs []string
	status  string
	time    int32
}

func NewMsgSender() service.IMessageSender {
	return &MsgSender{
		mutex:   &sync.RWMutex{},
		clients: make(map[string]*Client),
		rooms:   make(map[string]*Room),
		players: make(map[string]*resources.Player),
	}
}

// Send implements service.MessageSender.
func (s *MsgSender) Send(ctx context.Context, to string, data interface{}) error {
	s.mutex.RLock()
	client, ok := s.clients[to]
	s.mutex.RUnlock()
	if !ok {
		return errors.New("client not found")
	}
	fmt.Println("YYYYY")
	select {
	case client.ch <- data:
		return nil
	case <-time.After(writeWait):
		return errors.New("websocket write timeout")
	}
}

func (s *MsgSender) IsPlayerRegistered(playerID string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, ok := s.players[playerID]
	return ok
}

func (s *MsgSender) Broadcast(ctx context.Context, roomID string, data interface{}) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, id := range s.rooms[roomID].userIDs {
		client, ok := s.clients[id]
		if !ok {
			continue
		}

		select {
		case client.ch <- data:
		case <-time.After(writeWait):
			return errors.New("websocket write timeout")
		}

	}
	return nil
}

func (s *MsgSender) Register(roomID string, player *resources.Player, conn *websocket.Conn, err chan error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	client := &Client{
		conn:   conn,
		cancel: make(chan struct{}),
		ch:     make(chan interface{}, 100),
		err:    err,
	}
	go client.run()

	s.clients[player.PlayerId] = client
	room, exists := s.rooms[roomID]
	if !exists {
		room = &Room{
			userIDs: make([]string, 0, 4),
			time:    int32(-1*constants.CountDownTimer - 1),
		}
	}

	room.userIDs = append(room.userIDs, player.PlayerId)
	s.rooms[roomID] = room
	s.players[player.PlayerId] = player
}

func (s *MsgSender) Unregister(userID, roomID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	client, ok := s.clients[userID]
	if !ok {
		return
	}

	close(client.cancel)
	delete(s.clients, userID)
	delete(s.players, userID)

	room := s.rooms[roomID]
	for i, id := range room.userIDs {
		if id == userID {
			room.userIDs = append(room.userIDs[:i], room.userIDs[i+1:]...)
			s.rooms[roomID] = room
			break
		}
	}
}

func (s *MsgSender) UpdatePlayer(player *resources.Player) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.players[player.PlayerId]; !ok {
		return errors.New("player not found")
	}

	s.players[player.PlayerId] = player
	return nil
}

func (s *MsgSender) SetRoomStatus(roomID string, status string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	room := s.rooms[roomID]
	room.status = status
	s.rooms[roomID] = room
	return nil
}

func (s *MsgSender) GetRoomStatus(roomID string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if _, ok := s.rooms[roomID]; !ok {
		return "", errors.New("roomID not found")
	}
	return s.rooms[roomID].status, nil
}

func (s *MsgSender) GetPlayersInRoom(roomID string) ([]*resources.Player, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	room, exist := s.rooms[roomID]
	if !exist {
		return make([]*resources.Player, 0), nil
	}
	players := make([]*resources.Player, 0, len(room.userIDs))
	for _, playerID := range room.userIDs {
		if player, ok := s.players[playerID]; ok {
			players = append(players, player)
		}
	}
	return players, nil
}

func (s *MsgSender) GetTime(ctx context.Context, roomID string) int32 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	room, exist := s.rooms[roomID]
	if !exist {
		return 0
	}
	t := room.time
	if t > 0 {
		return int32(constants.TimeOutTimer-1) - t
	} else {
		return t
	}
}

func (s *MsgSender) setTime(ctx context.Context, roomID string, time int32) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	room, exist := s.rooms[roomID]
	if !exist {
		return
	}
	room.time = time
	s.rooms[roomID] = room
}

func (s *MsgSender) incrimentTime(ctx context.Context, roomID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	room, exist := s.rooms[roomID]
	if !exist {
		return
	}
	room.time++
	s.rooms[roomID] = room
}

func (s *MsgSender) StartTimer(ctx context.Context, roomID string) {
	ticker := time.NewTicker(time.Duration(constants.IntervalTicker) * time.Second)
	defer ticker.Stop()

	timer := time.NewTimer(time.Duration(constants.CountDownTimer+constants.TimeOutTimer) * time.Second)
	s.mutex.RLock()
	_, exist := s.rooms[roomID]
	s.mutex.RUnlock()
	if !exist {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p, err := s.GetPlayersInRoom(roomID)
			if err != nil {
				fmt.Println("err: ", err)
			}
			s.incrimentTime(ctx, roomID)
			r := &rpc.GameStatusResponse{
				RoomId:  roomID,
				Event:   resources.Event_EVENT_TIMER,
				Players: p,
				Time:    s.GetTime(ctx, roomID),
				Mode:    resources.Mode_MODE_MULTI,
			}
			fmt.Println("response: ", r)
			data, err := proto.Marshal(r)
			if err != nil {
				fmt.Println("err: ", err)
			}
			s.Broadcast(ctx, roomID, data)
		case <-timer.C:
			return
		}
	}
}

func (s *MsgSender) DestroyRoom(ctx context.Context, roomID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, exist := s.rooms[roomID]
	if !exist {
		return
	}
	for _, id := range s.rooms[roomID].userIDs {
		client, ok := s.clients[id]
		if !ok {
			continue
		}
		close(client.cancel)
		delete(s.clients, id)
		delete(s.players, id)
	}

	delete(s.rooms, roomID)

}
