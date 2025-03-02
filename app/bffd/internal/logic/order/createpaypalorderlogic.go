package order

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/order/rpc/order"
	"ymir.com/pkg/id"
	"ymir.com/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePaypalOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePaypalOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaypalOrderLogic {
	return &CreatePaypalOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePaypalOrderLogic) CreatePaypalOrder(req *types.CreatePaypalOrderRequest) (resp *types.CreatePaypalOrderResponse, err error) {
	uId, err := id.GetDecodedUserId(l.ctx)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.ErrorNotAuthorized)
	}

	oId, err := id.DecodeId(req.OrderId)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.ErrorReuqestParam)
	}

	respbCheck, err := l.svcCtx.OrderRPC.GetOrder(l.ctx, &order.GetOrderRequest{
		UserId:  uId,
		OrderId: oId,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "[CreatePaypalOrder] get order failed, orderId: %d", oId)
	}
	if respbCheck.Order == nil {
		return nil, xerr.NewErrCode(xerr.ErrorResourceForbiden)
	}

	respb, err := l.svcCtx.OrderRPC.CreatePaypalOrder(l.ctx, &order.CreatePaypalOrderRequest{
		RequestId: req.RequestId,
		OrderId:   oId,
		UserId:    uId,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "[CreatePaypalOrder] create paypal order failed, orderId: %d", oId)
	}

	return &types.CreatePaypalOrderResponse{
		OrderId:       req.OrderId,
		PaypalOrderId: respb.PaypalOrderId,
	}, nil
}
