package order

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/order/rpc/order"
	"ymir.com/app/product/rpc/product"
	"ymir.com/pkg/id"
	"ymir.com/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
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

	res, err := mr.MapReduce(func(source chan<- *order.OrderItem) {
		for i := 0; i < len(respb.OrderItems); i++ {
			source <- respb.OrderItems[i]
		}
	}, func(item *order.OrderItem, writer mr.Writer[*types.OrderItem], cancel func(error)) {
		var p *product.ProductDetailResponse
		var pc *product.ProductColorResponse
		err := mr.Finish(func() error {
			var err error
			p, err = l.svcCtx.ProductRPC.ProductDetail(l.ctx, &product.ProductDetailReqeust{
				Id: item.ProductId,
			})
			return err
		}, func() error {
			var err error
			pc, err = l.svcCtx.ProductRPC.ProductColor(l.ctx, &product.ProductColorRequest{
				ColorId: item.ColorId,
			})
			return err
		})
		if err != nil {
			cancel(err)
		}
		writer.Write(&types.OrderItem{
			ProductId:      id.EncodeId(p.Id),
			ProductColorId: id.EncodeId(pc.Color.Id),
			Size:           item.Size,
			Quantity:       item.Qunantity,
		})
	}, func(pipe <-chan *types.OrderItem, writer mr.Writer[[]types.OrderItem], cancel func(error)) {
		var items []types.OrderItem
		for item := range pipe {
			items = append(items, *item)
		}
		writer.Write(items)
	})
	if err != nil {
		return nil, errors.Wrapf(err, "get order items by orderId: %d failed", oId)
	}

	return &types.GetOrderItemResponse{
		OrderItems: res,
		Total:      int64(len(res)),
	}, nil
}
