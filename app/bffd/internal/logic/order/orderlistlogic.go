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

type OrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderListLogic {
	return &OrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderListLogic) OrderList(req *types.GetOrderRequest) (*types.GetOrderResponse, error) {
	if req.PageSize > 20 {
		return nil, xerr.NewErrCode(xerr.ErrorReuqestParam)
	}
	uId, err := id.GetDecodedUserId(l.ctx)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.ErrorNotAuthorized)
	}

	orders, err := l.svcCtx.OrderRPC.OrderList(l.ctx, &order.GetOrderListRequest{
		UserId:   uId,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	var res types.GetOrderResponse

	for i := range orders.Order {
		res.Orders = append(res.Orders, types.Order{
			OrderId: id.EncodeId(orders.Order[i].OrderId),
			Status:  orders.Order[i].Status,
			Price:   orders.Order[i].TotalPrice,
			Unit:    orders.Order[i].Unit,
		})
	}
	res.Total = orders.Total

	return &res, nil
}
