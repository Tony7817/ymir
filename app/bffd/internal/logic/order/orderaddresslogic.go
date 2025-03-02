package order

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/order/rpc/order"
	"ymir.com/pkg/id"
	"ymir.com/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderAddressLogic {
	return &OrderAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderAddressLogic) OrderAddress(req *types.GetOrderAddressRequest) (*types.GetOrderAddressResponse, error) {
	uId, err := id.GetDecodedUserId(l.ctx)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.ErrorNotAuthorized)
	}

	oId, err := id.DecodeId(req.OrderId)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.ErrorReuqestParam)
	}

	respb, err := l.svcCtx.OrderRPC.OrderAddress(l.ctx, &order.GetOrderAddressRequest{
		UserId:  uId,
		OrderId: oId,
	})
	if err != nil {
		return nil, err
	}

	if respb.Address == nil {
		return nil, xerr.NewErrCode(xerr.DataNoExistError)
	}

	return &types.GetOrderAddressResponse{
		Address: types.OrderAddress{
			AddressLine1: respb.Address.AddressLine1,
			AddressLine2: respb.Address.AddressLine2,
			AdminArea1:   respb.Address.AdminArea1,
			AdminArea2:   respb.Address.AdminArea2,
			PostalCode:   respb.Address.PostalCode,
			CountryCode:  respb.Address.CountryCode,
			Email:        respb.Address.Email,
			FullName:     respb.Address.FullName,
		},
	}, nil
}
