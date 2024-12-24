package logic

import (
	"context"

	"ymir.com/app/product/rpc/internal/svc"
	"ymir.com/app/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddProductAmountInCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddProductAmountInCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddProductAmountInCartLogic {
	return &AddProductAmountInCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddProductAmountInCartLogic) AddProductAmountInCart(in *product.AddProductAmountInCartRequest) (*product.AddProductAmountInCartResponse, error) {
	if err := l.svcCtx.ProductCartModel.IncrProductAmount(l.ctx, in.UserId, in.ProductId, in.Color); err != nil {
		return nil, err
	}

	return &product.AddProductAmountInCartResponse{}, nil
}
