package logic

import (
	"context"

	"ymir.com/app/product/admin/internal/svc"
	"ymir.com/app/product/admin/product"
	"ymir.com/app/product/model"

	"github.com/pkg/errors"
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
		offset = (in.Page - 1) * in.PageSize
		p      []*model.Product
		total  int64
	)
	err := mr.Finish(func() error {
		var err error
		p, err = l.svcCtx.ProductModel.FindProductList(l.ctx, in.StarId, offset, in.PageSize)
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
		return nil, errors.Wrapf(err, "[ProductList] failed to get product list")
	}

	productItems, err := l.colorOfProducts(p)
	if err != nil {
		return nil, err
	}

	var res = &product.ProductListResponse{
		Total:    total,
		Products: productItems,
	}

	return res, nil
}

func (l *ProductListLogic) colorOfProducts(ps []*model.Product) ([]*product.ProductListItem, error) {
	return mr.MapReduce(func(source chan<- *model.Product) {
		for _, p := range ps {
			source <- p
		}
	}, func(p *model.Product, writer mr.Writer[*product.ProductListItem], cancel func(error)) {
		color, err := l.svcCtx.ProductColorModel.FindOneByProductIdIsDefault(l.ctx, p.Id, 1)
		if err != nil {
			cancel(errors.Wrapf(err, "[colorOfProducts] failed to find color of product"))
			return
		}
		writer.Write(&product.ProductListItem{
			Id:          p.Id,
			Description: p.Description,
			Name:        p.Name,
			Rate:      p.Rate,
			RateCount: p.RateCount,
			DefaultColor: &product.ProductListColorItem{
				CoverUrl: color.CoverUrl,
				Price:    color.Price,
				Unit:     color.Unit,
			},
		})
	}, func(pipe <-chan *product.ProductListItem, writer mr.Writer[[]*product.ProductListItem], cancel func(error)) {
		var items []*product.ProductListItem
		for p := range pipe {
			items = append(items, p)
		}
		writer.Write(items)
	})
}
