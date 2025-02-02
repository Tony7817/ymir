package logic

import (
	"context"

	"ymir.com/app/order/rpc/internal/svc"
	"ymir.com/app/order/rpc/order"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type OrderItemsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderItemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderItemsLogic {
	return &OrderItemsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderItemsLogic) OrderItems(in *order.GetOrderItemsRequest) (*order.GetOrderItemsResponse, error) {
	ois, err := l.svcCtx.OrderItemModel.FindOneByOrderIdNotSoftDelete(l.ctx, in.OrderId)
	if err != nil {
		return nil, errors.Wrapf(err, "find order items by order id: %d", in.OrderId)
	}

	var orderItems []*order.OrderItem

	for i := 0; i < len(ois); i++ {
		orderItems = append(orderItems, &order.OrderItem{
			ProductId:   ois[i].ProductId,
			ColorId:     ois[i].ProductColorId,
			Size:        ois[i].Size,
			Qunantity:   ois[i].Quantity,
			Price:       ois[i].Price,
			OrderItemId: ois[i].Id,
		})
	}

	return &order.GetOrderItemsResponse{
		OrderItems: orderItems,
	}, nil
}
