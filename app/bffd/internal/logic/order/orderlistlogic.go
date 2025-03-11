package order

import (
	"context"
	"sort"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/order/rpc/order"
	"ymir.com/pkg/id"
	"ymir.com/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
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

func (l *OrderListLogic) OrderList(req *types.GetOrderListRequest) (*types.GetOrderListResponse, error) {
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

	type indexedOrder struct {
		index int
		order *order.OrderContent
	}
	var orderWithIndex = make([]indexedOrder, len(orders.Order))
	for i := range orders.Order {
		orderWithIndex[i] = indexedOrder{
			index: i,
			order: orders.Order[i],
		}
	}

	type indexedOrderListItem struct {
		index int
		order *types.OrderListItem
	}

	ois, err := mr.MapReduce(func(source chan<- *indexedOrder) {
		for i := range orderWithIndex {
			source <- &orderWithIndex[i]
		}
	}, func(orderWithIdx *indexedOrder, writer mr.Writer[*indexedOrderListItem], cancel func(error)) {
		o, err := l.svcCtx.OrderRPC.GetOrderItem(l.ctx, &order.GetOrderItemRequest{
			OrderId: orderWithIdx.order.OrderId,
			UserId:  uId,
		})
		if err != nil {
			cancel(err)
			return
		}
		var orderLen = min(3, len(o.OrderItems))
		var coverUrls = make([]string, orderLen)
		for i := range orderLen {
			coverUrls[i] = o.OrderItems[i].ProductColorCoverUrl
		}
		writer.Write(&indexedOrderListItem{
			index: orderWithIdx.index,
			order: &types.OrderListItem{
				CreatedAt:      orderWithIdx.order.CreatedAt,
				OrderId:        id.EncodeId(orderWithIdx.order.OrderId),
				CoverUrls:      coverUrls,
				Status:         orderWithIdx.order.Status,
				Price:          orderWithIdx.order.TotalPrice,
				Unit:           orderWithIdx.order.Unit,
				OrderItemTotal: o.TotalItems,
			},
		})
	}, func(pipe <-chan *indexedOrderListItem, writer mr.Writer[[]types.OrderListItem], cancel func(error)) {
		var orderItems []indexedOrderListItem
		for item := range pipe {
			orderItems = append(orderItems, *item)
		}
		sort.Slice(orderItems, func(i, j int) bool {
			return orderItems[i].index < orderItems[j].index
		})
		var orders = make([]types.OrderListItem, len(orderItems))
		for i := range orderItems {
			orders[i] = *orderItems[i].order
		}
		writer.Write(orders)
	})
	if err != nil {
		return nil, err
	}

	return &types.GetOrderListResponse{
		Orders: ois,
		Total:  orders.Total,
	}, nil
}
