package order

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/order/rpc/order"
	"ymir.com/app/product/rpc/product"
	"ymir.com/pkg/id"
	"ymir.com/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
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

	oItems, err := mr.MapReduce(func(source chan<- *order.OrderItem) {
		for i := range o.OrderItems {
			source <- o.OrderItems[i]
		}
	}, func(item *order.OrderItem, writer mr.Writer[*types.OrderItem], cancel func(error)) {
		var p *product.ProductDetailResponse
		var c *product.ProductColorResponse
		err := mr.Finish(func() error {
			var err error
			p, err = l.svcCtx.ProductRPC.ProductDetail(l.ctx, &product.ProductDetailReqeust{
				Id: item.ProductId,
			})
			if err != nil {
				return err
			}
			return nil
		}, func() error {
			var err error
			c, err = l.svcCtx.ProductRPC.ProductColor(l.ctx, &product.ProductColorRequest{
				ColorId: item.ColorId,
			})
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			cancel(err)
			return
		}

		oItem := &types.OrderItem{
			ProductId:          id.EncodeId(p.Id),
			ProductDescription: p.Description,
			ProductCoverImg:    c.Color.CoverUrl,
			ColorId:            id.EncodeId(c.Color.Id),
			Size:               item.Size,
			Quantity:           item.Qunantity,
			Price:              c.Color.Price,
		}
		writer.Write(oItem)
	}, func(pipe <-chan *types.OrderItem, writer mr.Writer[[]types.OrderItem], cancel func(error)) {
		var items []types.OrderItem
		for item := range pipe {
			items = append(items, *item)
		}
		writer.Write(items)
	})
	if err != nil {
		return nil, err
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
