package logic

import (
	"context"

	"ymir.com/app/order/model"
	"ymir.com/app/order/rpc/internal/svc"
	"ymir.com/app/order/rpc/order"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type GetOrderItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderItemLogic {
	return &GetOrderItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOrderItemLogic) GetOrderItem(in *order.GetOrderItemRequest) (*order.GetOrderItemResponse, error) {
	var orderItems []model.OrderItem
	var total int64
	err := mr.Finish(func() error {
		var err error
		orderItems, err = l.svcCtx.OrderItemModel.FindOrderItemsByOrderIdUserId(l.ctx, in.OrderId, in.UserId)
		if err != nil {
			return errors.Wrapf(err, "find order items by orderId: %d failed", in.OrderId)
		}
		return nil
	}, func() error {
		var err error
		total, err = l.svcCtx.OrderItemModel.CountOrderItemsByOrderIdUserId(l.ctx, in.OrderId, in.UserId)
		if err != nil {
			return errors.Wrapf(err, "count order items by orderId: %d failed", in.OrderId)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var res []*order.OrderItem
	for i := range orderItems {
		res = append(res, &order.OrderItem{
			ProductId:            orderItems[i].ProductId,
			ColorId:              orderItems[i].ProductColorId,
			Size:                 orderItems[i].Size,
			Qunantity:            orderItems[i].Quantity,
			Price:                orderItems[i].Price,
			OrderItemId:          orderItems[i].Id,
			ProductDescription:   orderItems[i].ProductDescription,
			Color:                orderItems[i].ProductColor,
			ProductColorCoverUrl: orderItems[i].ProductColorCoverUrl,
		})
	}

	return &order.GetOrderItemResponse{
		OrderItems: res,
		TotalItems:      total,
	}, nil
}
