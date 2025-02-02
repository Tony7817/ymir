package logic

import (
	"context"
	"database/sql"

	"ymir.com/app/order/model"
	"ymir.com/app/order/rpc/internal/svc"
	"ymir.com/app/order/rpc/order"
	"ymir.com/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type GetOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOrderLogic) GetOrder(in *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	var orderRes *model.Order
	var oiRes []*order.OrderItem
	err := mr.Finish(func() error {
		var err error
		orderRes, err = l.svcCtx.OrderModel.FindOne(l.ctx, in.OrderId)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errors.Wrapf(err, "find order by id: %d failed, user id: %d", in.OrderId, in.UserId)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Wrapf(xerr.NewErrCode(xerr.ErrorOrderNotExist), "cannot find order by orderId: %d, userId: %d", in.OrderId, in.UserId)
		}
		return nil
	}, func() error {
		var err error
		res, err := l.svcCtx.OrderItemModel.FindOneByOrderIdNotSoftDelete(l.ctx, in.OrderId)
		if err != nil {
			return errors.Wrapf(err, "find order items by orderId: %d failed", in.OrderId)
		}
		for i := 0; i < len(res); i++ {
			oiRes = append(oiRes, &order.OrderItem{
				ProductId:   res[i].ProductId,
				ColorId:     res[i].ProductColorId,
				Size:        res[i].Size,
				Qunantity:   res[i].Quantity,
				Price:       res[i].Price,
				OrderItemId: res[i].Id,
			})
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "get order by orderId: %d failed", in.OrderId)
	}

	return &order.GetOrderResponse{
		Order: &order.OrderContent{
			OrderId:    orderRes.Id,
			UserId:     orderRes.UserId,
			TotalPrice: orderRes.TotalPrice,
			Status:     orderRes.Status,
		},
		OrderItems: oiRes,
	}, nil
}
