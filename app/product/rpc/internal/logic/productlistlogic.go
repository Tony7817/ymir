package logic

import (
	"context"
	"sort"

	"ymir.com/app/product/model"
	"ymir.com/app/product/rpc/internal/svc"
	"ymir.com/app/product/rpc/product"

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
	for i := 0; i < len(p); i++ {
		l.Logger.Debugf("product: %+v", p[i].Description)
	}

	productItems, err := l.colorOfProducts(p)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(productItems); i++ {
		l.Logger.Debugf("productItem: %+v", productItems[i].Description)
	}

	var res = &product.ProductListResponse{
		Total:    total,
		Products: productItems,
	}

	return res, nil
}

func (l *ProductListLogic) colorOfProducts(ps []*model.Product) ([]*product.ProductListItem, error) {
	type indexedProduct struct {
		index   int
		product *model.Product
	}
	type indexedProductListItem struct {
		index int
		item  *product.ProductListItem
	}
	var productsWithIndex = make([]*indexedProduct, len(ps))

	for i := 0; i < len(ps); i++ {
		productsWithIndex[i] = &indexedProduct{
			index:   i,
			product: ps[i],
		}
	}
	res, err := mr.MapReduce(func(source chan<- *indexedProduct) {
		for _, p := range productsWithIndex {
			source <- p
		}
	}, func(p *indexedProduct, writer mr.Writer[*indexedProductListItem], cancel func(error)) {
		color, err := l.svcCtx.ProductColorModel.FindOneByProductIdIsDefault(l.ctx, p.product.Id, 1)
		if err != nil {
			cancel(errors.Wrapf(err, "[colorOfProducts] failed to find color of product"))
			return
		}
		writer.Write(&indexedProductListItem{
			index: p.index,
			item: &product.ProductListItem{
				Id:          p.product.Id,
				Description: p.product.Description,
				Unit:        color.Unit,
				DefaultColor: &product.ProductListColorItem{
					CoverUrl: color.CoverUrl,
					Price:    color.Price,
					Unit:     color.Unit,
				},
			},
		})
	}, func(pipe <-chan *indexedProductListItem, writer mr.Writer[[]*product.ProductListItem], cancel func(error)) {
		var items []*indexedProductListItem
		for p := range pipe {
			items = append(items, p)
		}
		sort.Slice(items, func(i, j int) bool {
			return items[i].index < items[j].index
		})
		var res = make([]*product.ProductListItem, len(items))
		for i := 0; i < len(items); i++ {
			res[i] = items[i].item
		}
		writer.Write(res)
	})
	if err != nil {
		return nil, errors.Wrapf(err, "[colorOfProducts] failed to map reduce")
	}

	return res, nil
}
