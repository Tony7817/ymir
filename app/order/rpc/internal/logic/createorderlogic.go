package logic

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ymir.com/app/order/model"
	"ymir.com/app/order/rpc/internal/svc"
	"ymir.com/app/order/rpc/order"
	"ymir.com/pkg/xerr"

	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmgrpc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrderLogic) CreateOrder(in *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	if len(in.Items) > 10 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ErrorRequestOrderMaximunReach), "Request too many orders, order num: %d", len(in.Items))
	}

	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		l.Logger.Errorf("[CreateOrder] BarrierFromGrpc error: %+v", err)
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	err = barrier.CallWithDB(l.svcCtx.DB, func(tx *sql.Tx) error {
		err := mr.Finish(func() error {
			var totalPrice int64
			for i := 0; i < len(in.Items); i++ {
				totalPrice += in.Items[i].Price * in.Items[i].Qunantity
			}
			_, err := l.svcCtx.OrderModel.SFInsertTx(l.ctx, tx, &model.Order{
				Id:         in.OrderId,
				UserId:     in.UserId,
				Status:     model.OrderStatusPending,
				TotalPrice: totalPrice,
			})
			if err != nil {
				l.Logger.Errorf("[CreateOrder] OrderModel.SFInsertTx error: %+v", err)
				return status.Error(codes.Aborted, dtmcli.ResultFailure)
			}
			return nil
		}, func() error {
			err := l.createOrderItems(tx, in.OrderId, in.Items)
			if err != nil {
				return status.Error(codes.Aborted, dtmcli.ResultFailure)
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &order.CreateOrderResponse{
		OrderId: in.OrderId,
	}, nil
}

func (l *CreateOrderLogic) createOrderItems(tx *sql.Tx, oId int64, os []*order.OrderItem) error {
	return mr.MapReduceVoid(func(source chan<- *order.OrderItem) {
		for i := 0; i < len(os); i++ {
			source <- os[i]
		}
	},
		func(o *order.OrderItem, writer mr.Writer[any], cancel func(error)) {
			_, err := l.svcCtx.OrderItemModel.SFInsertTx(l.ctx, tx, &model.OrderItem{
				Id:             o.OrderItemId,
				OrderId:        oId,
				ProductId:      o.ProductId,
				ProductColorId: o.ColorId,
				Size:           o.Size,
				Quantity:       o.Qunantity,
				Price:          o.Price,
				Subtotal:       o.Qunantity * o.Price,
			})
			if err != nil {
				l.Logger.Errorf("[CreateOrder] OrderItemModel.SFInsertTx error: %+v", err)
				cancel(status.Error(codes.Aborted, dtmcli.ResultFailure))
			}
		},
		func(pipe <-chan any, cancel func(error)) {})
}
