package logic

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ymir.com/app/product/rpc/internal/svc"
	"ymir.com/app/product/rpc/product"

	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type CheckoutProductRollbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckoutProductRollbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckoutProductRollbackLogic {
	return &CheckoutProductRollbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckoutProductRollbackLogic) CheckoutProductRollback(in *product.CheckoutProductRequest) (*product.CheckoutProductResponse, error) {
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	err = barrier.CallWithDB(l.svcCtx.DB, func(tx *sql.Tx) error {
		err := mr.MapReduceVoid(func(source chan<- *product.ProductCartItem) {
			for i := range in.Items {
				source <- in.Items[i]
			}
		}, func(item *product.ProductCartItem, writer mr.Writer[*product.ProductCartItem], cancel func(error)) {
			pc, err := l.svcCtx.ProductCartCheckoutModel.FindOneByOrderIdColorIdProductIdSize(l.ctx, in.Orderid, item.ColorId, item.ProductId, item.Size)
			if err != nil {
				cancel(err)
			}

			err = l.svcCtx.ProductCartModel.CheckoutProductRollback(l.ctx, tx, pc)
			if err != nil {
				cancel(err)
			}
		}, func(pipe <-chan *product.ProductCartItem, cancel func(error)) {
			for range pipe {
			}
		})
		if err != nil {
			return status.Error(codes.Aborted, dtmcli.ResultFailure)
		}
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	return &product.CheckoutProductResponse{
		Status: product.Status_OK,
	}, nil
}
