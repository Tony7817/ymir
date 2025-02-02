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
	err := l.svcCtx.ProductCartModel.InsertIntoProductCart(l.ctx, in.ProductCartId, in.UserId, in.ProductId, in.Size, in.ColorId)
	if err != nil {
		return nil, err
	}

	return &product.AddProductToCartResponse{
		ProductCartId: in.ProductCartId,
	}, nil
}
