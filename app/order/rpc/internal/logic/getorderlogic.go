package logic

import (
	"context"
	"database/sql"

	"ymir.com/app/order/rpc/internal/svc"
	"ymir.com/app/order/rpc/order"
	"ymir.com/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
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
	res, err := l.svcCtx.OrderModel.FindOne(l.ctx, in.OrderId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrapf(err, "find order by id: %d failed, user id: %d", in.OrderId, in.UserId)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ErrorOrderNotExist), "cannot find order by orderId: %d", in.OrderId)
	}

	return &order.GetOrderResponse{
		OrderId:    res.Id,
		UserId:     res.UserId,
		TotalPrice: res.TotalPrice,
		Status:     res.Status,
	}, nil
}
