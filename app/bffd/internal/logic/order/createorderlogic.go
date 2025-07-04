package order

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/order/model"
	"ymir.com/app/order/rpc/order"
	"ymir.com/app/product/rpc/product"
	"ymir.com/pkg/id"
	"ymir.com/pkg/vars"
	"ymir.com/pkg/xerr"

	"github.com/dtm-labs/dtmgrpc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/threading"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	sf     *id.Snowflake
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		sf:     id.NewSnowFlake(),
	}
}

func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderRequest) (*types.CreateOrderResponse, error) {
	uId, err := id.GetDecodedUserId(l.ctx)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.ErrorNotAuthorized)
	}
	orderRpcServer, err := l.svcCtx.Config.OrderRPC.BuildTarget()
	if err != nil {
		return nil, errors.Wrap(err, "[CreateOrder] get order target path failed")
	}
	productRpcServer, err := l.svcCtx.Config.ProductRPC.BuildTarget()
	if err != nil {
		return nil, errors.Wrap(err, "[CreateOrder] get product target path failed")
	}

	isOrderIdempotent, err := l.checkOrderIdempotent(req.RequestId)
	if err != nil {
		return nil, err
	}
	if !isOrderIdempotent {
		return nil, xerr.NewErrCode(xerr.ErrorIdempotence)
	}

	// prepare rpc request data
	var oitems = make([]*order.OrderItem, len(req.Orders))
	var psitems = make([]*product.ProductStockItem, len(req.Orders))
	var pcitems = make([]*product.ProductCartItem, len(req.Orders))
	for i := range req.Orders {
		pId, err := id.DecodeId(req.Orders[i].ProductId)
		if err != nil {
			return nil, errors.Wrap(err, "[CreateOrder] deocde pid failed")
		}
		cId, err := id.DecodeId(req.Orders[i].ProductColorId)
		if err != nil {
			return nil, errors.Wrap(err, "[CreateOrder] decode cid failed")
		}
		oiId := id.SF.GenerateID()
		oitems[i] = &order.OrderItem{
			ProductId:            pId,
			ColorId:              cId,
			Size:                 req.Orders[i].Size,
			Qunantity:            req.Orders[i].Quantity,
			OrderItemId:          oiId,
			ProductDescription:   req.Orders[i].ProductDescription,
			Color:                req.Orders[i].Color,
			ProductColorCoverUrl: req.Orders[i].ProductCoverUrl,
			Price:                0,
		}
		psitems[i] = &product.ProductStockItem{
			ProductId: pId,
			ColorId:   cId,
			Size:      req.Orders[i].Size,
			Quantity:  req.Orders[i].Quantity,
		}
		pcitems[i] = &product.ProductCartItem{
			ProductId: pId,
			ColorId:   cId,
			Size:      req.Orders[i].Size,
		}
	}

	// get the price
	oitems, err = l.checkOrderPrice(oitems)
	if err != nil {
		return nil, err
	}

	var gid = dtmgrpc.MustGenGid(vars.DtmServer)
	var oId = id.SF.GenerateID()
	saga := dtmgrpc.NewSagaGrpc(vars.DtmServer, gid).
		Add(orderRpcServer+order.Order_CreateOrder_FullMethodName, orderRpcServer+order.Order_CreateOrderRollback_FullMethodName, &order.CreateOrderRequest{
			OrderId:   oId,
			Items:     oitems,
			UserId:    uId,
			RequestId: req.RequestId,
		}).
		Add(productRpcServer+product.Product_DecreaseProductStockOfOrder_FullMethodName, productRpcServer+product.Product_IncreaseProductStockOfOrder_FullMethodName, &product.DecreaseProductStockRequest{
			ProductStockItem: psitems,
		}).
		Add(productRpcServer+product.Product_CheckoutProduct_FullMethodName, productRpcServer+product.Product_CheckoutProductRollback_FullMethodName, &product.CheckoutProductRequest{
			Orderid: oId,
			UserId:  uId,
			Items:   pcitems,
		}).EnableConcurrent()
	saga.WaitResult = true
	err = saga.Submit()
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ErrorCreateOrder), "[CreateOrder] dtm trans failed, err: %+v", err)
	}

	// cache request id
	threading.GoSafe(func() {
		_, _ = l.svcCtx.OrderRPC.CheckOrderIdempotent(l.ctx, &order.CheckOrderIdempotentRequest{
			RequestId: req.RequestId,
		})
	})

	return &types.CreateOrderResponse{
		OrderId: id.EncodeId(oId),
		Status:  model.OrderStatusPending,
	}, nil
}

func (l *CreateOrderLogic) checkOrderPrice(orders []*order.OrderItem) ([]*order.OrderItem, error) {
	res, err := mr.MapReduce(func(source chan<- *order.OrderItem) {
		for _, order := range orders {
			source <- order
		}
	}, func(item *order.OrderItem, writer mr.Writer[*order.OrderItem], cancel func(error)) {
		p, err := l.svcCtx.ProductRPC.ProductColor(l.ctx, &product.ProductColorRequest{
			ColorId: item.ColorId,
		})
		if err != nil {
			cancel(err)
			return
		}
		item.Price = p.Color.Price
		writer.Write(item)
	}, func(pipe <-chan *order.OrderItem, writer mr.Writer[[]*order.OrderItem], cancel func(error)) {
		var items []*order.OrderItem
		for item := range pipe {
			items = append(items, item)
		}
		writer.Write(items)
	})
	if err != nil {
		return nil, errors.Wrap(err, "[CreateOrder] checkOrderPrice failed")
	}

	return res, nil
}

func (l *CreateOrderLogic) checkOrderIdempotent(requestId string) (bool, error) {
	o, err := l.svcCtx.OrderRPC.CheckOrderIdempotent(l.ctx, &order.CheckOrderIdempotentRequest{
		RequestId: requestId,
	})
	if err != nil {
		return false, err
	}

	return o.IsIdempotent, nil
}
