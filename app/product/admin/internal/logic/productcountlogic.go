package logic

import (
	"context"

	"ymir.com/app/product/admin/internal/svc"
	"ymir.com/app/product/admin/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductCountLogic {
	return &ProductCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductCountLogic) ProductCount(in *product.ProductCountRequest) (*product.ProductCountResponse, error) {
	total, err := l.svcCtx.ProductModel.CountTotalProduct(l.ctx, in.StarId)
	if err != nil {
		return nil, err
	}

	return &product.ProductCountResponse{
		Count: total,
	}, nil
}
