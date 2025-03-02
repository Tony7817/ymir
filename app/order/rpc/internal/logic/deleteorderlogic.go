package logic

import (
	"context"

	"ymir.com/app/order/rpc/internal/svc"
	"ymir.com/app/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOrderLogic {
	return &DeleteOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteOrderLogic) DeleteOrder(in *order.DeleteOrderRequest) (*order.DeleteOrderResponse, error) {
	if err := l.svcCtx.OrderModel.Delete(l.ctx, in.OrderId); err != nil {
		return nil, err
	}

	return &order.DeleteOrderResponse{
		OrderId: in.OrderId,
	}, nil
}
