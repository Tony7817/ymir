package stock

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DecreaseStockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDecreaseStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DecreaseStockLogic {
	return &DecreaseStockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DecreaseStockLogic) DecreaseStock(req *types.DecreaseStockRequest) (resp *types.DecreaseStockResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
