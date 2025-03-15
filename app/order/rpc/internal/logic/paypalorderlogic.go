package logic

import (
	"context"
	"database/sql"
	"errors"

	"ymir.com/app/order/rpc/internal/svc"
	"ymir.com/app/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type PaypalOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPaypalOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaypalOrderLogic {
	return &PaypalOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PaypalOrderLogic) PaypalOrder(in *order.PaypalOrderReuqest) (*order.PaypalOrderResponse, error) {
	po, err := l.svcCtx.PaypalModel.FindOneByOrderIdUserId(l.ctx, in.OrderId, in.UserId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return &order.PaypalOrderResponse{
			Paypal: nil,
		}, nil
	}

	return &order.PaypalOrderResponse{
		Paypal: &order.Paypal{
			Id:            po.Id,
			OrderId:       po.OrderId,
			RequestId:     po.RequestId,
			PaypalOrderId: po.PaypalOrderId,
		},
	}, nil
}
