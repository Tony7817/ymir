package logic

import (
	"context"
	"strings"

	"ymir.com/app/product/rpc/internal/svc"
	"ymir.com/app/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductColorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductColorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductColorLogic {
	return &ProductColorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductColorLogic) ProductColor(in *product.ProductColorRequest) (*product.ProductColorResponse, error) {
	color, err := l.svcCtx.ProductColorModel.FindOne(l.ctx, in.ColorId)
	if err != nil {
		return nil, err
	}

	return &product.ProductColorResponse{
		Id:             color.Id,
		ProductId:      color.ProductId,
		Color:          color.Color,
		Images:         strings.Split(color.Images, ","),
		DetailImages:   strings.Split(color.DetailImages, ","),
		Price:          color.Price,
		Unit:           color.Unit,
		AvaliableSizes: strings.Split(color.AvailableSize, ","),
	}, nil
}
