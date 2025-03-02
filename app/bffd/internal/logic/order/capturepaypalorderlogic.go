package order

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/order/model"
	"ymir.com/app/order/rpc/order"
	"ymir.com/pkg/id"
	"ymir.com/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CapturePaypalOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCapturePaypalOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CapturePaypalOrderLogic {
	return &CapturePaypalOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CapturePaypalOrderLogic) CapturePaypalOrder(req *types.CapturePaypalOrderRequest) (resp *types.CapturePaypalOrderResponse, err error) {
	uId, err := id.GetDecodedUserId(l.ctx)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.ErrorNotAuthorized)
	}

	oId, err := id.DecodeId(req.OrderId)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.ErrorReuqestParam)
	}

	o, err := l.svcCtx.OrderRPC.GetOrder(l.ctx, &order.GetOrderRequest{
		UserId:  uId,
		OrderId: oId,
	})
	if err != nil {
		return nil, err
	}
	if o.Order == nil {
		return nil, xerr.NewErrCode(xerr.ErrorResourceForbiden)
	}

	respb, err := l.svcCtx.OrderRPC.CaptureOrder(l.ctx, &order.CapturePaypalOrderRequest{
		RequestId: req.OrderId,
		UserId:    uId,
		OrderId:   oId,
	})
	if err != nil {
		return nil, err
	}

	if respb.Status == "COMPLETED" {
		// change order status
		_, err = l.svcCtx.OrderRPC.UpdateOrder(l.ctx, &order.UpdateOrderRequest{
			OrderId: oId,
			UserId:  uId,
			Status:  &model.OrderStatusPaid,
		})
		if err != nil {
			return nil, err
		}
	}

	return &types.CapturePaypalOrderResponse{
		OrderId: id.EncodeId(respb.OrderId),
		Status:  respb.Status,
		PayerInfo: types.PayerInfo{
			Email:         respb.PayerInfo.Email,
			PayerFullName: respb.PayerInfo.PayerFullName,
			PayerAddress: types.PayerAddress{
				AddressLine1: respb.PayerInfo.PayerAddress.AddressLine1,
				AddressLine2: respb.PayerInfo.PayerAddress.AddressLine2,
				AdminArea1:   respb.PayerInfo.PayerAddress.AdminArea1,
				AdminArea2:   respb.PayerInfo.PayerAddress.AdminArea2,
				PostalCode:   respb.PayerInfo.PayerAddress.PostalCode,
				CountryCode:  respb.PayerInfo.PayerAddress.CountryCode,
			},
		},
	}, nil
}
