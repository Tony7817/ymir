package logic

import (
	"context"
	"encoding/json"

	"ymir.com/app/product/rpc/internal/svc"
	"ymir.com/app/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductListLogic {
	return &ProductListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductListLogic) ProductList(in *product.ProductListRequest) (*product.ProductListResponse, error) {
	var offset = (in.Page - 1) * in.PageSize
	ps, err := l.svcCtx.ProductModel.FindProductList(l.ctx, in.StarId, offset, in.PageSize)
	logx.Infof("%+v", ps)
	if err != nil {
		return nil, err
	}

	var productItems []*product.ProductListItem
	for i := 0; i < len(ps); i++ {
		var images []*product.ProductImage
		if err := json.Unmarshal([]byte(ps[i].Images), &images); err != nil {
			return nil, err
		}
		productItems = append(productItems, &product.ProductListItem{
			Id:          ps[i].Id,
			Description: ps[i].Description,
			Color:       ps[i].Color,
			Coverurl:    images[0],
			Price:       ps[i].Price,
			Unit:        ps[i].Unit,
		})
	}

	var res = &product.ProductListResponse{
		Total:    int64(len(productItems)),
		Products: productItems,
	}

	return res, nil
}
