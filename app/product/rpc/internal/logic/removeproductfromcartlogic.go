package logic

import (
	"context"

	"ymir.com/app/product/rpc/internal/svc"
	"ymir.com/app/product/rpc/product"

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
	err := l.svcCtx.ProductCartModel.DeleteProductFromCart(l.ctx, in.ProductId, in.ColorId, in.UserId)
	if err != nil {
		return nil, err
	}

	return &product.RemoveProductFromCartResponse{}, nil
}
