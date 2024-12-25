package logic

import (
	"context"

	"ymir.com/app/product/rpc/internal/svc"
	"ymir.com/app/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductStockLogic {
	return &ProductStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductStockLogic) ProductStock(in *product.ProductStockRequest) (*product.ProductStockResponse, error) {
	s, err := l.svcCtx.ProductStockModel.FindOneByProductIdColorIdSize(l.ctx, in.ProductId, in.ColorId, in.Size)
	if err != nil {
		return nil, err
	}

	return &product.ProductStockResponse{
		Stock: s.InStock,
	}, nil
}
