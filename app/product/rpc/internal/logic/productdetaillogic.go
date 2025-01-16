package logic

import (
	"context"

	"ymir.com/app/product/rpc/internal/svc"
	"ymir.com/app/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductDetailLogic {
	return &ProductDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductDetailLogic) ProductDetail(in *product.ProductDetailReqeust) (*product.ProductDetailResponse, error) {
	p, err := l.svcCtx.ProductModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	var res = &product.ProductDetailResponse{
		Id:             p.Id,
		Description:    p.Description,
		Rate:           p.Rate,
		ReteCount:      p.RateCount,
		SoldNum:        p.SoldNum,
		Detail:         p.Detail.String,
		StarId:         p.StarId,
	}

	return res, nil
}
