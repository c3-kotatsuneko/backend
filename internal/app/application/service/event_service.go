package service

import (
	"context"

	"github.com/c3-kotatsuneko/backend/internal/domain/service"
)

type IEventService interface {
	TimeAttack(ctx context.Context, senderID, roomID string) error
	MultiPlay(ctx context.Context, senderID, roomID string) error
	Training(ctx context.Context, senderID, roomID string) error
}

type EventService struct {
	msgSender service.IMessageSender
}

func NewEventService(msgSender service.IMessageSender) IEventService {
	return &EventService{
		msgSender: msgSender,
	}
}

func (s *EventService) TimeAttack(ctx context.Context, senderID, roomID string) error {
	return nil
}

func (s *EventService) MultiPlay(ctx context.Context, senderID, roomID string) error {
	return nil
}

func (s *EventService) Training(ctx context.Context, senderID, roomID string) error {
	return nil
}
