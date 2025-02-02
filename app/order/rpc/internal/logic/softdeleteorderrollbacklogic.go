package logic

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ymir.com/app/order/rpc/internal/svc"
	"ymir.com/app/order/rpc/order"

	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SoftDeleteOrderRollbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSoftDeleteOrderRollbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SoftDeleteOrderRollbackLogic {
	return &SoftDeleteOrderRollbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SoftDeleteOrderRollbackLogic) SoftDeleteOrderRollback(in *order.SoftDeleteOrderRequest) (*order.SoftDeleteOrderResponse, error) {
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	err = barrier.CallWithDB(l.svcCtx.DB, func(tx *sql.Tx) error {
		return l.svcCtx.OrderModel.SoftDeleteRollbackTx(l.ctx, tx, in.OrderId, in.UserId)
	})
	if err != nil {
		l.Logger.Errorf("[SoftDeleteOrder] SoftDeleteTx error: %+v", err)
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	return &order.SoftDeleteOrderResponse{
		OrderId: in.OrderId,
	}, nil
}
