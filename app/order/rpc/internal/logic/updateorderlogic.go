package logic

import (
	"context"

	"ymir.com/app/order/model"
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
	var o model.Order
	if in.Status != nil {
		o.Status = *in.Status
	}

	if err := l.svcCtx.OrderModel.Update(l.ctx, &o); err != nil {
		return nil, err
	}

	return &order.UpdateOrderResponse{
		Status: o.Status,
	}, nil
}
