package logic

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ymir.com/app/product/rpc/internal/svc"
	"ymir.com/app/product/rpc/product"
	"ymir.com/pkg/id"

	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmgrpc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type CheckoutProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckoutProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckoutProductLogic {
	return &CheckoutProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckoutProductLogic) CheckoutProduct(in *product.CheckoutProductRequest) (*product.CheckoutProductResponse, error) {
	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	err = barrier.CallWithDB(l.svcCtx.DB, func(tx *sql.Tx) error {
		err := mr.MapReduceVoid(func(source chan<- *product.ProductCartItem) {
			for i := range in.Items {
				source <- in.Items[i]
			}
		}, func(item *product.ProductCartItem, writer mr.Writer[int64], cancel func(error)) {
			pc, err := l.svcCtx.ProductCartModel.FindOneByUserIdColorIdProductIdSize(l.ctx, in.UserId, item.ColorId, item.ProductId, item.Size)
			if err != nil {
				l.Logger.Errorf("[CheckoutProduct] ProductCartModel.FindOneByUserIdColorIdProductIdSize error: %+v", err)
				cancel(err)
			}
			l.Logger.Debugf("[CheckoutProduct] ProductCartModel.FindOneByUserIdColorIdProductIdSize pc: %+v", pc)
			err = l.svcCtx.ProductCartModel.CheckoutProduct(l.ctx, tx, in.Orderid, id.SF.GenerateID(), pc)
			if err != nil {
				l.Logger.Errorf("[CheckoutProduct] ProductCartModel.CheckoutProduct error: %+v", err)
				cancel(err)
			}
		}, func(pipe <-chan int64, cancel func(error)) {
			for range pipe {
			}
		})
		if err != nil {
			return status.Error(codes.Aborted, dtmcli.ResultFailure)
		}

		return nil
	})
	if err != nil {
		l.Logger.Errorf("[CheckoutProduct] barrier.CallWithDB error: %+v", err)
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	return &product.CheckoutProductResponse{
		Status: product.Status_OK,
	}, nil
}
