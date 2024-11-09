package service

import (
	"github.com/c3-kotatsuneko/backend/internal/cat/physics"
	"github.com/c3-kotatsuneko/backend/internal/cat/repository"
	"github.com/c3-kotatsuneko/backend/internal/domain/entity"
)

type IObjectService interface {
	GetObjectsSlice() []*entity.Nekojarashi
	GetObjectByObjID(key string) *entity.Nekojarashi
	CollideWithObj(roomID string) map[string][]string
	ApplyForceToObj(obj1ID, obj2ID string)
	InitObjects(roomID string)
	UpdatePosition(roomID string)
	SharePosition(roomID string, objs []*entity.Nekojarashi)
}

type ObjectService struct {
	or repository.IObjectRepository
	nr repository.INikukyuRepository
	cr repository.ICatHouseRepository
}

func NewObjectService(or repository.IObjectRepository, nr repository.INikukyuRepository, cr repository.ICatHouseRepository) IObjectService {
	return &ObjectService{
		or: or,
		nr: nr,
		cr: cr,
	}
}

func (os *ObjectService) GetObjectsSlice() []*entity.Nekojarashi {
	return os.or.GetObjectsSlice()
}

func (os *ObjectService) GetObjectByObjID(key string) *entity.Nekojarashi {
	return os.or.GetObjectByObjID(key)
}

// roomの全オブジェクトの衝突判定
func (os *ObjectService) CollideWithObj(roomID string) map[string][]string {
	catHouse := os.cr.GetCatHouseByRoomID(roomID)
	allObj := os.or.GetObjectsByObjIDs(catHouse.Nekojarashis)
	collidedObjIDs := make(map[string][]string, len(allObj))
	for i := 0; i < len(allObj); i++ {
		for j := 0; j <= i; j++ {
			if allObj[i].ID == allObj[j].ID {
				continue
			}
			if collided := physics.IsColliding(allObj[i].Position, allObj[j].Position); collided {
				collidedObjIDs[allObj[i].ID] = append(collidedObjIDs[allObj[i].ID], allObj[j].ID)
			}
		}
	}

	return collidedObjIDs
}

func (os *ObjectService) ApplyForceToObj(obj1ID, obj2ID string) {
	obj1 := os.or.GetObjectByObjID(obj1ID)
	obj2 := os.or.GetObjectByObjID(obj2ID)
	physics.CollidedVelocity(obj1, obj2)
	// physics.UpdatePosition(obj1)
	// physics.UpdatePosition(obj2)
}

func (os *ObjectService) InitObjects(roomID string) {
	catHouse := os.cr.GetCatHouseByRoomID(roomID)
	for _, v := range catHouse.Nekojarashis {
		os.or.InitObjects(v)
	}
}

func (os *ObjectService) UpdatePosition(roomID string) {
	catHouse := os.cr.GetCatHouseByRoomID(roomID)
	for _, v := range catHouse.Nekojarashis {
		obj := os.or.GetObjectByObjID(v)
		physics.UpdateVelocity(obj)
		physics.UpdatePosition(obj)
	}
}

func (os *ObjectService) SharePosition(roomID string, objs []*entity.Nekojarashi) {
	for _, obj := range objs {
		os.or.ModifyObjects(obj.ID, obj)
	}
}
