package logic

import (
	"context"

	"ymir.com/app/order/rpc/internal/svc"
	"ymir.com/app/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrderLogic {
	return &UpdateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateOrderLogic) UpdateOrder(in *order.UpdateOrderRequest) (*order.UpdateOrderResponse, error) {
	o, err := l.svcCtx.OrderModel.FindOne(l.ctx, in.OrderId)
	if err != nil {
		return nil, err
	}
	o.Status = *in.Status

	if err := l.svcCtx.OrderModel.Update(l.ctx, o); err != nil {
		return nil, err
	}

	return &order.UpdateOrderResponse{
		Status: o.Status,
	}, nil
}
