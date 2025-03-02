package logic

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ymir.com/app/product/rpc/internal/svc"
	"ymir.com/app/product/rpc/product"
	"ymir.com/pkg/xerr"

	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmgrpc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type DecreaseProductStockOfOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDecreaseProductStockOfOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DecreaseProductStockOfOrderLogic {
	return &DecreaseProductStockOfOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DecreaseProductStockOfOrderLogic) DecreaseProductStockOfOrder(in *product.DecreaseProductStockRequest) (*product.DecreaseProductStockResponse, error) {
	if len(in.ProductStockItem) > 10 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ErrorRequestOrderMaximunReach), "Request too many product stocks, order num: %d", len(in.ProductStockItem))
	}

	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		l.Logger.Errorf("[DecreaseProductStockOfOrder] BarrierFromGrpc error: %v", err)
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	err = barrier.CallWithDB(l.svcCtx.DB, func(tx *sql.Tx) error {
		err := l.decreaseStocks(tx, in.ProductStockItem)
		if err != nil {
			l.Logger.Errorf("[DecreaseProductStockOfOrder] decreaseStocks error: %v", err)
			return status.Error(codes.Aborted, dtmcli.ResultFailure)
		}
		return nil
	})
	if err != nil {
		l.Logger.Errorf("[DecreaseProductStockOfOrder] decreaseStocks error: %v", err)
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	return &product.DecreaseProductStockResponse{
		Status: product.Status_OK,
	}, nil
}

func (l *DecreaseProductStockOfOrderLogic) decreaseStocks(tx *sql.Tx, s []*product.ProductStockItem) error {
	return mr.MapReduceVoid(func(source chan<- *product.ProductStockItem) {
		for i := 0; i < len(s); i++ {
			source <- s[i]
		}
	}, func(s *product.ProductStockItem, writer mr.Writer[any], cancel func(error)) {
		err := l.svcCtx.ProductStockModel.DecreaseProductStockTx(l.ctx, tx, s.ProductId, s.ColorId, s.Quantity, s.Size)
		l.Logger.Infof("%+v", err)
		if err != nil {
			l.Logger.Errorf("[DecreaseProductStockOfOrder] ProductStockModel.DecreaseProductStockTx error: %+v", err)
			cancel(status.Error(codes.Aborted, dtmcli.ResultFailure))
		}
	}, func(pipe <-chan any, cancel func(error)) {})
}
