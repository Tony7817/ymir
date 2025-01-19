package stock

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IncreaseStockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIncreaseStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IncreaseStockLogic {
	return &IncreaseStockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IncreaseStockLogic) IncreaseStock(req *types.IncreaseStockRequest) (resp *types.IncreaseStockResponse, err error) {

	return
}
