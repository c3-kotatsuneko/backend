package repository

import (
	"sync"

	"github.com/c3-kotatsuneko/backend/internal/cat/constants"
	"github.com/c3-kotatsuneko/backend/internal/domain/entity"
	"github.com/c3-kotatsuneko/backend/pkg/uuid"
)

type ICatHouseRepository interface {
	InitCatHouse(roomID string)
	GetCatHouseByRoomID(roomID string) *entity.CatHouse
	AddNikukyu(roomID, userID string)
}

type CatHouseRepository struct {
	mux       sync.RWMutex
	catHouses map[string]*entity.CatHouse
}

func NewCatHouseRepository() ICatHouseRepository {
	return &CatHouseRepository{
		mux:       sync.RWMutex{},
		catHouses: make(map[string]*entity.CatHouse),
	}
}

func (cr *CatHouseRepository) GetCatHouseByRoomID(roomID string) *entity.CatHouse {
	cr.mux.RLock()
	defer cr.mux.RUnlock()
	if _, ok := cr.catHouses[roomID]; !ok {
		return nil
	}
	return cr.catHouses[roomID].DeepCopy()
}

func (cr *CatHouseRepository) InitCatHouse(roomID string) {
	cr.mux.Lock()
	defer cr.mux.Unlock()
	nekojarashi := uuid.NewUUIDs(constants.InitBlock)
	nikukyu := make([]string, 0, 4)
	cr.catHouses[roomID] = &entity.CatHouse{
		Nekojarashis: nekojarashi,
		Nikukyus:     nikukyu,
	}
}

func (cr *CatHouseRepository) AddNikukyu(roomID, userID string) {
	cr.mux.Lock()
	defer cr.mux.Unlock()
	cr.catHouses[roomID].Nikukyus = append(cr.catHouses[roomID].Nikukyus, userID)
}
