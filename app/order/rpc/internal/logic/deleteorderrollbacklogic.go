package logic

import (
	"context"

	"ymir.com/app/order/rpc/internal/svc"
	"ymir.com/app/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteOrderRollbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteOrderRollbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOrderRollbackLogic {
	return &DeleteOrderRollbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteOrderRollbackLogic) DeleteOrderRollback(in *order.DeleteOrderRequest) (*order.DeleteOrderResponse, error) {
	// todo: add your logic here and delete this line

	return &order.DeleteOrderResponse{}, nil
}
