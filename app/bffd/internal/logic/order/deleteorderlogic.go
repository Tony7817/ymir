package order

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/order/rpc/order"
	"ymir.com/app/product/rpc/product"
	"ymir.com/pkg/id"
	"ymir.com/pkg/vars"
	"ymir.com/pkg/xerr"

	"github.com/dtm-labs/dtmgrpc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOrderLogic {
	return &DeleteOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteOrderLogic) DeleteOrder(req *types.DeleteOrderRequest) (*types.DeleteORderResponse, error) {
	uId, err := id.GetDecodedUserId(l.ctx)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.ErrorNotAuthorized)
	}

	oId, err := id.DecodeId(req.OrderId)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.ErrorReuqestParam)
	}

	orders, err := l.svcCtx.OrderRPC.GetOrder(l.ctx, &order.GetOrderRequest{
		UserId:  uId,
		OrderId: oId,
	})
	if err != nil {
		return nil, err
	}

	var psItems []*product.ProductStockItem
	for i := 0; i < len(orders.OrderItems); i++ {
		psItems = append(psItems, &product.ProductStockItem{
			ProductId: orders.OrderItems[i].ProductId,
			ColorId:   orders.OrderItems[i].ColorId,
			Size:      orders.OrderItems[i].Size,
			Quantity:  orders.OrderItems[i].Qunantity,
		})
	}

	var gid = dtmgrpc.MustGenGid(vars.DtmServer)
	orderRpcServer, err := l.svcCtx.Config.OrderRPC.BuildTarget()
	if err != nil {
		return nil, errors.Wrapf(err, "[DeleteOrder] get order target path failed")
	}
	productRpcServer, err := l.svcCtx.Config.ProductRPC.BuildTarget()
	if err != nil {
		return nil, errors.Wrapf(err, "[DeleteOrder] get product target path failed")
	}
	saga := dtmgrpc.NewSagaGrpc(vars.DtmServer, gid).
		Add(orderRpcServer+order.Order_SoftDeleteOrder_FullMethodName, orderRpcServer+order.Order_SoftDeleteOrderRollback_FullMethodName, &order.SoftDeleteOrderRequest{
			UserId:  uId,
			OrderId: oId,
		}).
		Add(productRpcServer+product.Product_IncreaseProductStockOfOrder_FullMethodName, productRpcServer+product.Product_DecreaseProductStockOfOrder_FullMethodName, &product.IncreaseProductStockRequest{
			ProductStockItem: psItems,
		}).EnableConcurrent()
	saga.WaitResult = true
	err = saga.Submit()
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ErrorDeleteOrder), "[DeleteOrder] submit saga failed, err: %+v", err)
	}

	return &types.DeleteORderResponse{
		OrderId: id.EncodeId(oId),
	}, nil
}
