package service

import (
	"github.com/c3-kotatsuneko/backend/internal/cat/constants"
	"github.com/c3-kotatsuneko/backend/internal/cat/physics"
	"github.com/c3-kotatsuneko/backend/internal/cat/repository"
	"github.com/c3-kotatsuneko/backend/internal/domain/entity"
)

type IHandService interface {
	CalculateHandForce(hand *entity.Nikukyu) *entity.Vector3    // 手の力の更新
	CollideWithObj(roomID string, hand *entity.Nikukyu) *string // ブロックとの衝突判定(衝突したブロックのIDを返す)
	TransferHandToNikukyu(hand *entity.Hand) *entity.Nikukyu
	ApplyForceToObj(id string, force *entity.Vector3) // ブロックに力を加える
}

type HandService struct {
	or repository.IObjectRepository
	nr repository.INikukyuRepository
	cr repository.ICatHouseRepository
}

func NewHand(or repository.IObjectRepository, nr repository.INikukyuRepository, cr repository.ICatHouseRepository) IHandService {
	return &HandService{
		or: or,
		nr: nr,
		cr: cr,
	}
}

func (h *HandService) CalculateHandForce(hand *entity.Nikukyu) *entity.Vector3 {
	currHandVel := physics.CalculateVelocity(hand.PrevActionPosition, hand.ActionPosition, constants.TimeStep)
	handAcc := physics.CalculateAcceleration(hand.PrevVelocity, *currHandVel, constants.TimeStep)
	handForce := physics.CalculateForce(constants.BlockMass, *handAcc)
	return handForce
}

// roomの全Objectとの当たり判定を行い、当たったObjectのIDを返す
func (h *HandService) CollideWithObj(roomID string, hand *entity.Nikukyu) *string {
	catHouse := h.cr.GetCatHouseByRoomID(roomID)
	for _, v := range catHouse.Nekojarashis {
		obj := h.or.GetObjectByObjID(v)
		if collided := physics.IsColliding(hand.ActionPosition, obj.Position); collided {
			return &obj.ID
		}
	}
	return nil
}

func (h *HandService) TransferHandToNikukyu(hand *entity.Hand) *entity.Nikukyu {
	nikukyu := h.nr.TransferHandToNikukyu(hand)
	return nikukyu
}

func (h *HandService) ApplyForceToObj(id string, force *entity.Vector3) {
	obj := h.or.GetObjectByObjID(id)
	physics.ApplyForce(obj, force)
	// physics.ApplyFriction(obj)
	// physics.UpdatePosition(obj)
}

func (h *HandService) InitNikukyu(roomID, userID string) {
	h.cr.AddNikukyu(roomID, userID)
}
