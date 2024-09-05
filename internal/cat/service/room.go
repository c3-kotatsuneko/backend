package service

import (
	"github.com/c3-kotatsuneko/backend/internal/cat/repository"
	"github.com/c3-kotatsuneko/backend/internal/domain/entity"
)

type IRoomService interface {
	InitRoom(roomID string)
	GetCatHouseByRoomID(roomID string) *entity.CatHouse
}

type RoomService struct {
	cr repository.ICatHouseRepository
}

func NewRoomService(cr repository.ICatHouseRepository) IRoomService {
	return &RoomService{
		cr: cr,
	}
}

func (r *RoomService) GetCatHouseByRoomID(roomID string) *entity.CatHouse {
	return r.cr.GetCatHouseByRoomID(roomID)
}

func (r *RoomService) InitRoom(roomID string) {
	r.cr.InitCatHouse(roomID)
}
