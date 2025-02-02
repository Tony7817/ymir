package logic

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ymir.com/app/order/rpc/internal/svc"
	"ymir.com/app/order/rpc/order"
	"ymir.com/pkg/xerr"

	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmgrpc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type CreateOrderRollbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderRollbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderRollbackLogic {
	return &CreateOrderRollbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrderRollbackLogic) CreateOrderRollback(in *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	l.Logger.Errorf("[CreateOrderRollback] CreateOrderRollback request: %+v", in)
	if len(in.Items) > 10 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ErrorRequestOrderMaximunReach), "Request too many orders, order num: %d", len(in.Items))
	}

	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}
	err = barrier.CallWithDB(l.svcCtx.DB, func(tx *sql.Tx) error {
		return mr.Finish(func() error {
			err := l.svcCtx.OrderModel.SoftDeleteTx(l.ctx, tx, in.OrderId, in.UserId)
			if err != nil {
				l.Logger.Errorf("[CreateOrderRollback] OrderModel.DeleteTx error: %+v", err)
				return status.Error(codes.Aborted, dtmcli.ResultFailure)
			}
			return nil
		}, func() error {
			return l.deleteOrderItmes(tx, in.Items)
		})
	})
	if err != nil {
		l.Logger.Errorf("[SoftDeleteOrder] SoftDeleteTx error: %+v", err)
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	return &order.CreateOrderResponse{
		OrderId: in.OrderId,
	}, nil
}

func (l *CreateOrderRollbackLogic) deleteOrderItmes(tx *sql.Tx, os []*order.OrderItem) error {
	return mr.MapReduceVoid(func(source chan<- int64) {
		for i := 0; i < len(os); i++ {
			source <- os[i].OrderItemId
		}
	}, func(id int64, writer mr.Writer[any], cancel func(error)) {
		err := l.svcCtx.OrderItemModel.SoftDeleteByOrderIdTx(l.ctx, tx, id)
		if err != nil {
			l.Logger.Errorf("[CreateOrderRollback] OrderItemModel.DeleteTx error: %+v", err)
			cancel(status.Error(codes.Aborted, dtmcli.ResultFailure))
		}
	}, func(pipe <-chan any, cancel func(error)) {})
}
