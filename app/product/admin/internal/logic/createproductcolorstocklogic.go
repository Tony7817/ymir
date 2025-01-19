package logic

import (
	"context"

	"ymir.com/app/product/admin/internal/svc"
	"ymir.com/app/product/admin/product"
	"ymir.com/app/product/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductColorStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateProductColorStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductColorStockLogic {
	return &CreateProductColorStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateProductColorStockLogic) CreateProductColorStock(in *product.CreateProductColorStockRequest) (*product.CreateProductColorStockResponse, error) {
	id, err := l.svcCtx.ProductStockModel.InsertSF(l.ctx, &model.ProductStock{
		Id:        in.ProductColorStockId,
		ProductId: in.ProductId,
		ColorId:   in.ColorId,
		InStock:   in.InStock,
		Size:      in.Size,
	})
	if err != nil {
		return nil, err
	}

	return &product.CreateProductColorStockResponse{
		Id: id,
	}, nil
}
