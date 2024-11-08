package constants

import "github.com/c3-kotatsuneko/protobuf/gen/game/resources"

const (
	IntervalTicker int = 1
	CountDownTimer int = 3
	//TODO: 本番では300にする
	TimeOutTimer      int    = 300
	StackBlock               = 5
	RoomStatusWaiting string = "waiting"
	RoomStatusPlaying string = "playing"
)

var Directions = []resources.Direction{
	resources.Direction_DIRECTION_FRONT,
	resources.Direction_DIRECTION_LEFT,
	resources.Direction_DIRECTION_BACK,
	resources.Direction_DIRECTION_RIGHT,
}
