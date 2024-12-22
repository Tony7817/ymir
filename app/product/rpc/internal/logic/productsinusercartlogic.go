package logic

import (
	"context"

	"ymir.com/app/product/rpc/internal/svc"
	"ymir.com/app/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductsInUserCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductsInUserCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductsInUserCartLogic {
	return &ProductsInUserCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductsInUserCartLogic) ProductsInUserCart(in *product.UserProductIdsRequest) (*product.UserProductIdsResponse, error) {
	

	return &product.UserProductIdsResponse{}, nil
}
