package logic

import (
	"context"

	"ymir.com/app/product/admin/internal/svc"
	"ymir.com/app/product/admin/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProductLogic {
	return &DeleteProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteProductLogic) DeleteProduct(in *product.DeleteProductRequest) (*product.DeleteProductResponse, error) {
	err := l.svcCtx.ProductModel.DeleteProduct(l.ctx, in.ProductId)
	if err != nil {
		return nil, err
	}

	return &product.DeleteProductResponse{
		ProductId: in.ProductId,
	}, nil
}
