package service

import (
	"context"

	"github.com/c3-kotatsuneko/backend/internal/domain/entity"
)

type ICatService interface {
	Do(ctx context.Context, roomID string, hand *entity.Hand) error        // 1フレームの間に行う処理
	Get(ctx context.Context, roomID string) ([]*entity.Object, error)      // Objectの座標を返す
	Init(ctx context.Context, roomID string) error                         // 初期化処理
	Share(ctx context.Context, roomID string, objs []*entity.Object) error // Objectの座標を共有する
}
