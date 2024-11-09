package cat

import (
	"context"

	"github.com/c3-kotatsuneko/backend/internal/cat/service"
	"github.com/c3-kotatsuneko/backend/internal/domain/entity"
	domainService "github.com/c3-kotatsuneko/backend/internal/domain/service"
)

type Cat struct {
	hs service.IHandService
	os service.IObjectService
	rs service.IRoomService
}

func New(hs service.IHandService, os service.IObjectService, rs service.IRoomService) domainService.ICatService {
	return &Cat{
		hs: hs,
		os: os,
		rs: rs,
	}
}

// 1フレームの間に行う処理
func (c *Cat) Do(ctx context.Context, roomID string, hand *entity.Hand) error {
	nikukyu := c.hs.TransferHandToNikukyu(hand)
	// 手の当たり判定を求める
	if id := c.hs.CollideWithObj(roomID, nikukyu); id != nil {
		// オブジェクトに与える影響を計算
		force := c.hs.CalculateHandForce(nikukyu)
		c.hs.ApplyForceToObj(*id, force)
	}
	// オブジェクト同士の当たり判定を求める
	for sourceObjID, targetObjID := range c.os.CollideWithObj(roomID) {
		// オブジェクトに与える影響を計算
		for _, id := range targetObjID {
			c.os.ApplyForceToObj(sourceObjID, id)
		}
	}
	// 全てのオブジェクトの位置を更新
	c.os.UpdatePosition(roomID)

	return nil
}

func (c *Cat) Get(ctx context.Context, roomID string) ([]*entity.Object, error) {
	catHouse := c.rs.GetCatHouseByRoomID(roomID)
	obj := make([]*entity.Object, 0, 20)
	for _, v := range catHouse.Nekojarashis {
		nekojarashi := c.os.GetObjectByObjID(v)
		obj = append(obj, &entity.Object{
			ID:       nekojarashi.ID,
			Layer:    nekojarashi.Layer,
			Kinds:    nekojarashi.Kinds,
			State:    nekojarashi.State,
			Position: nekojarashi.Position,
			Size:     nekojarashi.Size,
		})
	}
	return obj, nil
}

func (c *Cat) Init(ctx context.Context, roomID string) error {
	c.rs.InitRoom(roomID)
	c.os.InitObjects(roomID)

	return nil
}

func (c *Cat) Share(ctx context.Context, roomID string, objs []*entity.Object) error {
	catHouse := c.rs.GetCatHouseByRoomID(roomID)
	nekojarashi := make([]*entity.Nekojarashi, 0, 20)
	for _, obj := range objs {
		catHouse.Nekojarashis = append(catHouse.Nekojarashis, obj.ID)
		nekojarashi = append(nekojarashi, &entity.Nekojarashi{
			ID:       obj.ID,
			Layer:    obj.Layer,
			Kinds:    obj.Kinds,
			State:    obj.State,
			Position: obj.Position,
			Size:     obj.Size,
		})
	}
	c.os.SharePosition(roomID, nekojarashi)

	return nil
}
