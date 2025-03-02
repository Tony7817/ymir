package logic

import (
	"context"
	"database/sql"
	"errors"

	"ymir.com/app/order/rpc/internal/svc"
	"ymir.com/app/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderAddressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderAddressLogic {
	return &OrderAddressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderAddressLogic) OrderAddress(in *order.GetOrderAddressRequest) (*order.GetOrderAddressResponse, error) {
	oa, err := l.svcCtx.OrderAddressModel.FindOneByOrderId(l.ctx, in.OrderId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if errors.Is(err, sql.ErrNoRows) {
		return &order.GetOrderAddressResponse{
			Address: nil,
		}, nil
	}

	return &order.GetOrderAddressResponse{
		Address: &order.OrderAddress{
			AddressLine1: oa.Address1.String,
			AddressLine2: oa.Address2.String,
			AdminArea1:   oa.AdminArea1.String,
			AdminArea2:   oa.AdminArea2.String,
			PostalCode:   oa.PostalCode,
			CountryCode:  oa.PostalCode,
			Email:        oa.Email,
			FullName:     oa.FullName,
		},
	}, nil
}
