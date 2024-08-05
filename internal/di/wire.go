//go:build wireinject
// +build wireinject

package di

import (
	"github.com/c3-kotatsuneko/backend/internal/app/application/service"
	"github.com/c3-kotatsuneko/backend/internal/app/infrastructure"
	"github.com/c3-kotatsuneko/backend/internal/app/presentation"
	"github.com/c3-kotatsuneko/backend/internal/app/presentation/handler"
	"github.com/c3-kotatsuneko/backend/internal/app/presentation/switcher"
	"github.com/c3-kotatsuneko/backend/internal/app/presentation/websocket"
	"github.com/c3-kotatsuneko/backend/internal/cat"
	catRepository "github.com/c3-kotatsuneko/backend/internal/cat/repository"
	catService "github.com/c3-kotatsuneko/backend/internal/cat/service"

	"github.com/c3-kotatsuneko/backend/pkg/cache"
	"github.com/google/wire"
)

func InitHandler() *presentation.Root {
	wire.Build(
		cache.NewCacheClient,
		infrastructure.NewMsgSender,
		infrastructure.NewRoomObjectRepository,
		infrastructure.NewCat,
		service.NewRoomObjectService,
		service.NewEventService,
		switcher.NewPhysicsSwitcher,
		switcher.NewEventSwitcher,
		websocket.NewWSHandler,
		handler.NewPhysicsHandler,
		handler.NewEventHandler,
		presentation.New,

		cat.New,
		catService.NewHand,
		catService.NewObjectService,
		catRepository.NewHandRepository,
		catRepository.NewObjectRepository,
	)
	return &presentation.Root{}
}
