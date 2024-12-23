package logic

import (
	"context"

	"ymir.com/app/product/rpc/internal/svc"
	"ymir.com/app/product/rpc/model"
	"ymir.com/app/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
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
	var (
		offset   = (in.Page - 1) * in.PageSize
		products []*model.Product
		total    int64
	)
	err := mr.Finish(func() error {
		var err error
		products, err = l.svcCtx.ProductModel.FindProductList(l.ctx, in.StarId, offset, in.PageSize)
		if err != nil {
			return err
		}
		return nil
	}, func() error {
		var err error
		total, err = l.svcCtx.ProductModel.CountTotalProduct(l.ctx, in.StarId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var productItems []*product.ProductListItem
	for i := 0; i < len(products); i++ {
		productItems = append(productItems, &product.ProductListItem{
			Id:          productItems[i].Id,
			Description: productItems[i].Description,
			Color:       productItems[i].Color,
			Coverurl:    productItems[i].Coverurl,
			Price:       productItems[i].Price,
			Unit:        productItems[i].Unit,
		})
	}

	var res = &product.ProductListResponse{
		Total:    total,
		Products: productItems,
	}

	return res, nil
}
