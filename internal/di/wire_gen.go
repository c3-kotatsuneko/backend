// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	service2 "github.com/c3-kotatsuneko/backend/internal/app/application/service"
	"github.com/c3-kotatsuneko/backend/internal/app/infrastructure"
	"github.com/c3-kotatsuneko/backend/internal/app/presentation"
	"github.com/c3-kotatsuneko/backend/internal/app/presentation/handler"
	"github.com/c3-kotatsuneko/backend/internal/app/presentation/switcher"
	"github.com/c3-kotatsuneko/backend/internal/app/presentation/websocket"
	"github.com/c3-kotatsuneko/backend/internal/cat"
	"github.com/c3-kotatsuneko/backend/internal/cat/repository"
	"github.com/c3-kotatsuneko/backend/internal/cat/service"
	"github.com/c3-kotatsuneko/backend/pkg/cache"
)

// Injectors from wire.go:

func InitHandler() *presentation.Root {
	client := cache.NewCacheClient()
	iRoomObjectRepository := infrastructure.NewRoomObjectRepository(client)
	iMessageSender := infrastructure.NewMsgSender()
	iObjectRepository := repository.NewObjectRepository()
	iNikukyuRepository := repository.NewHandRepository()
	iCatHouseRepository := repository.NewCatHouseRepository()
	iHandService := service.NewHand(iObjectRepository, iNikukyuRepository, iCatHouseRepository)
	iObjectService := service.NewObjectService(iObjectRepository, iNikukyuRepository, iCatHouseRepository)
	iRoomService := service.NewRoomService(iCatHouseRepository)
	iCatService := cat.New(iHandService, iObjectService, iRoomService)
	iCatRepository := infrastructure.NewCat(iCatService)
	iRoomObjectService := service2.NewRoomObjectService(iRoomObjectRepository, iMessageSender, iCatRepository)
	iwsHandler := websocket.NewWSHandler(iRoomObjectService, iMessageSender)
	iPhysicsSwitcher := switcher.NewPhysicsSwitcher(iRoomObjectService, iMessageSender)
	physicsHandler := handler.NewPhysicsHandler(iwsHandler, iPhysicsSwitcher)
	iEventService := service2.NewEventService(iMessageSender)
	iEventSwitcher := switcher.NewEventSwitcher(iEventService, iMessageSender)
	eventHandler := handler.NewEventHandler(iwsHandler, iEventSwitcher)
	root := presentation.New(physicsHandler, eventHandler)
	return root
}
