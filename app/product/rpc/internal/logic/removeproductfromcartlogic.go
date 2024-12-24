package logic

import (
	"context"

	"ymir.com/app/product/rpc/internal/svc"
	"ymir.com/app/product/rpc/product"
	"ymir.com/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveProductFromCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveProductFromCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveProductFromCartLogic {
	return &RemoveProductFromCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveProductFromCartLogic) RemoveProductFromCart(in *product.RemoveProductFromCartRequest) (*product.RemoveProductFromCartResponse, error) {
	ret, err := l.svcCtx.ProductCartModel.SoftRemoveToProductFromCart(l.ctx, in.UserId, in.ProductId, in.Color)
	if err != nil {
		return nil, err
	}

	rowsAff, err := ret.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAff == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ServerCommonError), "product not found in cart")
	}

	return &product.RemoveProductFromCartResponse{}, nil
}
