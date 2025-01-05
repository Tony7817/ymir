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
	// todo: add your logic here and delete this line

	return &product.ProductCountResponse{}, nil
}
