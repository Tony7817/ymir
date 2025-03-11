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

type OrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailLogic {
	return &OrderDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderDetailLogic) OrderDetail(req *types.GetOrderDetailRequest) (resp *types.GetOrderDetailResponse, err error) {
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
		return nil, xerr.NewErrCode(xerr.DataNoExistError)
	}

	var oItems []types.OrderItem
	for i := range o.OrderItems {
		oItems = append(oItems, types.OrderItem{
			ProductId:          id.EncodeId(o.OrderItems[i].ProductId),
			ProductDescription: o.OrderItems[i].ProductDescription,
			ProductCoverUrl:    o.OrderItems[i].ProductColorCoverUrl,
			ProductColorId:     id.EncodeId(o.OrderItems[i].ColorId),
			Color:              o.OrderItems[i].Color,
			Size:               o.OrderItems[i].Size,
			Quantity:           o.OrderItems[i].Qunantity,
			Price:              o.OrderItems[i].Price,
		})
	}

	var res = types.GetOrderDetailResponse{
		Order: types.Order{
			RequestId: o.Order.RequestId,
			OrderId:   id.EncodeId(o.Order.OrderId),
			Status:    o.Order.Status,
			Price:     o.Order.TotalPrice,
			Unit:      o.Order.Unit,
		},
		OrderItems: oItems,
	}

	return &res, nil
}
