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

type IncreaseProductStockOfOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIncreaseProductStockOfOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IncreaseProductStockOfOrderLogic {
	return &IncreaseProductStockOfOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IncreaseProductStockOfOrderLogic) IncreaseProductStockOfOrder(in *product.DecreaseProductStockRequest) (*product.DecreaseProductStockResponse, error) {
	if len(in.ProductStockItem) > 10 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ErrorRequestOrderMaximunReach), "Request too many product stocks, order num: %d", len(in.ProductStockItem))
	}

	l.Logger.Debugf("increase stock")
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		l.Logger.Errorf("[IncreaseProductStockOfOrder] BarrierFromGrpc error: %v", err)
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	err = barrier.CallWithDB(l.svcCtx.DB, func(tx *sql.Tx) error {
		return l.increaseStocks(tx, in.ProductStockItem)
	})
	if err != nil {
		l.Logger.Errorf("[IncreaseProductStockOfOrder] increaseStocks error: %v", err)
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	return &product.DecreaseProductStockResponse{
		Status: product.Status_OK,
	}, nil
}

func (l *IncreaseProductStockOfOrderLogic) increaseStocks(tx *sql.Tx, s []*product.ProductStockItem) error {
	return mr.MapReduceVoid(func(source chan<- *product.ProductStockItem) {
		for i := 0; i < len(s); i++ {
			source <- s[i]
		}
	}, func(item *product.ProductStockItem, writer mr.Writer[any], cancel func(error)) {
		if err := l.svcCtx.ProductStockModel.IncreaseProductStockTx(l.ctx, tx, item.ProductId, item.ColorId, item.Quantity, item.Size); err != nil {
			l.Logger.Errorf("[IncreaseProductStockOfOrder] IncreaseProductStockTx error: %+v", err)
			cancel(status.Error(codes.Aborted, dtmcli.ResultFailure))
		}
	}, func(pipe <-chan any, cancel func(error)) {})
}
