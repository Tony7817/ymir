package logic

import (
	"context"
	"database/sql"
	"sort"

	"ymir.com/app/product/model"
	"ymir.com/app/product/rpc/internal/svc"
	"ymir.com/app/product/rpc/product"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type ProductsInCartListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductsInCartListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductsInCartListLogic {
	return &ProductsInCartListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductsInCartListLogic) ProductsInCartList(in *product.ProductsInCartListRequest) (*product.ProducrtsInCartListResponse, error) {
	productCartsInfo, err := l.svcCtx.ProductCartModel.FindProductsOfUser(in.UserId, in.Page, in.PageSize)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		productCartsInfo = []*model.ProductCart{}
	}

	var (
		products = make([]*product.ProductsInCartListItem, 0)
		total    int64
	)
	err = mr.Finish(
		func() error {
			var err error
			total, err = l.svcCtx.ProductCartModel.CountTotalProductOfUser(l.ctx, in.UserId)
			if err != nil {
				return err
			}
			return nil
		},
		func() error {
			var err error
			products, err = l.productsInCarts(productCartsInfo)
			if err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &product.ProducrtsInCartListResponse{
		Products: products,
		Total:    total,
	}, nil
}

func (l *ProductsInCartListLogic) productsInCarts(pcarts []*model.ProductCart) ([]*product.ProductsInCartListItem, error) {
	type indexedProductCart struct {
		index int
		item  *model.ProductCart
	}
	type indexedProductCartItem struct {
		index int
		item  *product.ProductsInCartListItem
	}
	var productCartsWithIndex = make([]*indexedProductCart, len(pcarts))
	for i := 0; i < len(pcarts); i++ {
		productCartsWithIndex[i] = &indexedProductCart{
			index: i,
			item:  pcarts[i],
		}
	}
	res, err := mr.MapReduce(func(source chan<- *indexedProductCart) {
		for _, item := range productCartsWithIndex {
			source <- item
		}
	}, func(pcart *indexedProductCart, writer mr.Writer[*indexedProductCartItem], cancel func(error)) {
		var p *model.Product
		var c *model.ProductColorDetail
		var s *model.ProductStock
		err := mr.Finish(func() error {
			var err error
			p, c, err = l.findProductDetailAndColorDetail(pcart.item.ProductId, pcart.item.ColorId)
			if err != nil {
				return err
			}
			return nil
		}, func() error {
			var err error
			s, err = l.svcCtx.ProductStockModel.FindOneByProductIdColorIdSize(l.ctx, pcart.item.ProductId, pcart.item.ColorId, pcart.item.Size)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			cancel(err)
			return
		}
		writer.Write(&indexedProductCartItem{
			index: pcart.index,
			item: &product.ProductsInCartListItem{
				ProductCartId: pcart.item.Id,
				ProductId:     p.Id,
				StarId:        p.StarId,
				ColorId:       c.Id,
				Amount:        pcart.item.Amount,
				Sizes:         pcart.item.Size,
				Description:   p.Description,
				Price:         c.Price,
				Unit:          c.Unit,
				CoverUrl:      c.CoverUrl,
				Stock:         s.InStock,
			},
		})
	}, func(pipe <-chan *indexedProductCartItem, writer mr.Writer[[]*product.ProductsInCartListItem], cancel func(error)) {
		productsItemWithIndex := []*indexedProductCartItem{}
		for item := range pipe {
			productsItemWithIndex = append(productsItemWithIndex, item)
		}
		sort.Slice(productsItemWithIndex, func(i, j int) bool {
			return productsItemWithIndex[i].index < productsItemWithIndex[j].index
		})

		res := make([]*product.ProductsInCartListItem, len(productsItemWithIndex))
		for i := 0; i < len(productsItemWithIndex); i++ {
			res[i] = productsItemWithIndex[i].item
		}
		writer.Write(res)
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get product detail and color detail")
	}

	return res, nil
}

func (l *ProductsInCartListLogic) findProductDetailAndColorDetail(productId int64, colorId int64) (*model.Product, *model.ProductColorDetail, error) {
	var c *model.ProductColorDetail
	var p *model.Product

	err := mr.Finish(func() error {
		var err error
		c, err = l.svcCtx.ProductColorModel.FindOne(l.ctx, colorId)
		if err != nil {
			return err
		}
		return nil
	}, func() error {
		var err error
		p, err = l.svcCtx.ProductModel.FindOne(l.ctx, productId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return p, c, nil
}
