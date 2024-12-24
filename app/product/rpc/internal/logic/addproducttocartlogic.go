package logic

import (
	"context"

	"ymir.com/app/product/rpc/internal/svc"
	"ymir.com/app/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddProductToCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddProductToCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddProductToCartLogic {
	return &AddProductToCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddProductToCartLogic) AddProductToCart(in *product.AddProductToCartRequest) (*product.AddProductToCartResponse, error) {
	ret, err := l.svcCtx.ProductCartModel.AddToProductCart(l.ctx, in.UserId, in.ProductId, in.Size, in.Color)
	if err != nil {
		return nil, err
	}

	lastInsertId, err := ret.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &product.AddProductToCartResponse{
		ProductCartId: lastInsertId,
	}, nil
}
