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

type OrderItemsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderItemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderItemsLogic {
	return &OrderItemsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderItemsLogic) OrderItems(req *types.GetOrderItemRequest) (*types.GetOrderItemResponse, error) {
	uId, err := id.GetDecodedUserId(l.ctx)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.ErrorNotAuthorized)
	}

	oId, err := id.DecodeId(req.OrderId)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.ErrorReuqestParam)
	}

	respb, err := l.svcCtx.OrderRPC.GetOrder(l.ctx, &order.GetOrderRequest{
		UserId:  uId,
		OrderId: oId,
	})
	if err != nil {
		return nil, err
	}

	var res types.GetOrderItemResponse

	for i := 0; i < len(respb.OrderItems); i++ {
		res.OrderItems = append(res.OrderItems, types.OrderItem{
			ProductId: id.EncodeId(respb.OrderItems[i].ProductId),
			ColorId:   id.EncodeId(respb.OrderItems[i].ColorId),
			Size:      respb.OrderItems[i].Size,
			Quantity:  respb.OrderItems[i].Qunantity,
			Price:     respb.OrderItems[i].Price,
		})
	}

	return &res, nil
}
