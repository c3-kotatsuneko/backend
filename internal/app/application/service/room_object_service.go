package service

import (
	"context"
	"fmt"

	"github.com/c3-kotatsuneko/backend/internal/domain/entity"
	"github.com/c3-kotatsuneko/backend/internal/domain/repository"
	"github.com/c3-kotatsuneko/backend/internal/domain/service"
)

type IRoomObjectService interface {
	Calculate(ctx context.Context, senderID, roomID string, hand *entity.Hand) error
	Get(ctx context.Context, roomID string) ([]*entity.Object, error)
	Init(ctx context.Context, roomID string) error
}

type RoomObjectService struct {
	roomObjRepo repository.IRoomObjectRepository
	msgSender   service.IMessageSender
	catRepo     repository.ICatRepository
}

func NewRoomObjectService(roomObjRepo repository.IRoomObjectRepository, msgSender service.IMessageSender, catRepo repository.ICatRepository) IRoomObjectService {
	return &RoomObjectService{
		roomObjRepo: roomObjRepo,
		msgSender:   msgSender,
		catRepo:     catRepo,
	}
}

func (s *RoomObjectService) Calculate(ctx context.Context, senderID, roomID string, hand *entity.Hand) error {
	fmt.Println("before calculate")

	if err := s.catRepo.Calculate(ctx, roomID, hand); err != nil {
		return err
	}

	return nil
}

func (s *RoomObjectService) Get(ctx context.Context, roomID string) ([]*entity.Object, error) {
	return s.catRepo.Get(ctx, roomID)
}

func (s *RoomObjectService) Init(ctx context.Context, roomID string) error {
	return s.catRepo.Init(ctx, roomID)
}
