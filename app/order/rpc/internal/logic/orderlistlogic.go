package logic

import (
	"context"

	"ymir.com/app/order/rpc/internal/svc"
	"ymir.com/app/order/rpc/order"
	"ymir.com/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type OrderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderListLogic {
	return &OrderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderListLogic) OrderList(in *order.GetOrderListRequest) (*order.GetOrderListResponse, error) {
	if in.PageSize > 20 {
		return nil, xerr.NewErrCode(xerr.ErrorReuqestParam)
	}

	var offset = (in.Page - 1) * in.PageSize
	var res order.GetOrderListResponse
	err := mr.Finish(func() error {
		orders, err := l.svcCtx.OrderModel.FindOrderList(l.ctx, in.UserId, offset, in.PageSize)
		if err != nil {
			return err
		}

		for i := range orders {
			res.Order = append(res.Order, &order.OrderContent{
				OrderId:    orders[i].Id,
				UserId:     orders[i].UserId,
				TotalPrice: orders[i].TotalPrice,
				Status:     orders[i].Status,
				Unit:       orders[i].Unit,
			})
		}

		return nil
	}, func() error {
		count, err := l.svcCtx.OrderModel.CountOrderList(l.ctx, in.UserId)
		if err != nil {
			return err
		}

		res.Total = count

		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "[OrderListLogic] OrderList failed")
	}

	return &res, nil
}
