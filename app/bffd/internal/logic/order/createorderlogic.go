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

	isOrderCreated, err := l.checkOrderIdempotent(req.RequestId)
	if err != nil {
		return nil, err
	}
	if isOrderCreated {
		return nil, xerr.NewErrCode(xerr.ErrorIdempotence)
	}

	var oitems []*order.OrderItem
	var psitems []*product.ProductStockItem
	for i := 0; i < len(req.Orders); i++ {
		pId, err := id.DecodeId(req.Orders[i].ProductId)
		if err != nil {
			return nil, errors.Wrap(err, "[CreateOrder] deocde pid failed")
		}
		cId, err := id.DecodeId(req.Orders[i].ColorId)
		if err != nil {
			return nil, errors.Wrap(err, "[CreateOrder] decode cid failed")
		}
		oiId := id.SF.GenerateID()
		oitems = append(oitems, &order.OrderItem{
			ProductId:   pId,
			ColorId:     cId,
			Size:        req.Orders[i].Size,
			Qunantity:   req.Orders[i].Quantity,
			OrderItemId: oiId,
		})
		psitems = append(psitems, &product.ProductStockItem{
			ProductId: pId,
			ColorId:   cId,
			Size:      req.Orders[i].Size,
			Quantity:  req.Orders[i].Quantity,
		})
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
		}).EnableConcurrent()
	saga.WaitResult = true
	err = saga.Submit()
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ErrorCreateOrder), "[CreateOrder] dtm trans failed, err: %+v", err)
	}

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
	respb, err := l.svcCtx.OrderRPC.PaypalOrder(l.ctx, &order.PaypalOrderReuqest{
		RequestId: requestId,
	})
	if err != nil {
		return false, err
	}

	if respb.Paypal == nil {
		return false, nil
	}

	return true, nil
}
