package logic

import (
	"context"

	"ymir.com/app/product/admin/internal/svc"
	"ymir.com/app/product/admin/product"
	"ymir.com/app/product/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductColorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateProductColorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductColorLogic {
	return &CreateProductColorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateProductColorLogic) CreateProductColor(in *product.CreateProductColorRequeset) (*product.CreateProductColorResponse, error) {
	id, err := l.svcCtx.ProductColorModel.SFInsert(l.ctx, &model.ProductColorDetail{
		Id:            in.ProductColorId,
		ProductId:     in.ProductId,
		Color:         in.ColorName,
		Images:        in.ImageUrl,
		DetailImages:  in.DetailImagesUrl,
		Price:         in.Price,
		Unit:          in.Unit,
		AvailableSize: in.Size,
		CoverUrl:      in.CoverUrl,
	})
	if err != nil {
		return nil, err
	}

	return &product.CreateProductColorResponse{
		Id: id,
	}, nil
}
