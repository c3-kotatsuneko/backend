package infrastructure

import (
	"context"

	"github.com/c3-kotatsuneko/backend/internal/domain/entity"
	"github.com/c3-kotatsuneko/backend/internal/domain/repository"
	domainService "github.com/c3-kotatsuneko/backend/internal/domain/service"
)

type Cat struct {
	cat domainService.ICatService
	// world map[string][]*entity.Object
	// mux   sync.RWMutex
}

func NewCat(cat domainService.ICatService) repository.ICatRepository {
	return &Cat{
		cat: cat,
		// world: make(map[string][]*entity.Object),
		// mux:   sync.RWMutex{},
	}
}

func (c *Cat) Calculate(ctx context.Context, roomID string, hand *entity.Hand) error {
	return c.cat.Do(ctx, roomID, hand)
}

func (c *Cat) Get(ctx context.Context, roomID string) ([]*entity.Object, error) {
	return c.cat.Get(ctx, roomID)
}

func (c *Cat) Init(ctx context.Context, roomID string) error {
	return c.cat.Init(ctx, roomID)
}

func (c *Cat) Share(ctx context.Context, roomID string, objs []*entity.Object) error {
	return c.cat.Share(ctx, roomID, objs)
}

// func (c *Cat) set(roomID string, objects []*entity.Object) error {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
// 	c.world[roomID] = append(c.world[roomID], objects...)
// 	return nil
// }

// func (c *Cat) get(roomID string) ([]*entity.Object, error) {
// 	c.mux.RLock()
// 	defer c.mux.RUnlock()
// 	originalSlice := c.world[roomID]
// 	slice := make([]*entity.Object, len(originalSlice))
// 	for i, obj := range originalSlice {
// 		slice[i] = obj.DeepCopy()
// 	}
// 	return slice, nil
// }
