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

type SoftDeleteOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSoftDeleteOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SoftDeleteOrderLogic {
	return &SoftDeleteOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SoftDeleteOrderLogic) SoftDeleteOrder(in *order.SoftDeleteOrderRequest) (*order.SoftDeleteOrderResponse, error) {
	if len(in.Items) > 10 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ErrorRequestOrderMaximunReach), "Request too many orders, order num: %d", len(in.Items))
	}

	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		//!!!一般数据库不会错误不需要dtm回滚，就让他一直重试，这时候就不要返回codes.Aborted, dtmcli.ResultFailure 就可以了，具体自己把控!!!
		return nil, status.Error(codes.Internal, err.Error())
	}
	err = barrier.CallWithDB(l.svcCtx.DB, func(tx *sql.Tx) error {
		return mr.Finish(func() error {
			return l.svcCtx.OrderModel.SoftDeleteTx(l.ctx, tx, in.OrderId, in.UserId)
		}, func() error {
			return l.deleteOrderItmes(tx, in.Items)
		})
	})
	if err != nil {
		l.Logger.Errorf("[SoftDeleteOrder] SoftDeleteTx error: %+v", err)
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	return &order.SoftDeleteOrderResponse{}, nil
}

func (l *SoftDeleteOrderLogic) deleteOrderItmes(tx *sql.Tx, os []*order.OrderItem) error {
	return mr.MapReduceVoid(func(source chan<- int64) {
		for i := 0; i < len(os); i++ {
			source <- os[i].OrderItemId
		}
	}, func(id int64, writer mr.Writer[any], cancel func(error)) {
		err := l.svcCtx.OrderItemModel.SoftDeleteByOrderIdTx(l.ctx, tx, id)
		if err != nil {
			cancel(err)
		}
	}, func(pipe <-chan any, cancel func(error)) {})
}
