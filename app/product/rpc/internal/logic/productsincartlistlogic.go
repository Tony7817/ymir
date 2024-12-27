package logic

import (
	"context"
	"database/sql"

	"ymir.com/app/product/rpc/internal/svc"
	"ymir.com/app/product/rpc/model"
	"ymir.com/app/product/rpc/product"

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
	return mr.MapReduce(func(source chan<- *model.ProductCart) {
		for _, item := range pcarts {
			source <- item
		}
	}, func(pcart *model.ProductCart, writer mr.Writer[*product.ProductsInCartListItem], cancel func(error)) {
		p, c, err := l.findProductDetailAndColorDetail(pcart.ProductId, pcart.ColorId)
		if err != nil {
			cancel(err)
			return
		}
		writer.Write(&product.ProductsInCartListItem{
			ProductId:   p.Id,
			StarId:      p.StarId,
			ColorId:     c.Id,
			Amount:      pcart.Amount,
			Sizes:       pcart.Size,
			Description: p.Description,
			Price:       c.Price,
			Unit:        c.Unit,
			CoverUrl:    c.CoverUrl,
		})
	}, func(pipe <-chan *product.ProductsInCartListItem, writer mr.Writer[[]*product.ProductsInCartListItem], cancel func(error)) {
		pitems := []*product.ProductsInCartListItem{}
		for item := range pipe {
			pitems = append(pitems, item)
		}
		writer.Write(pitems)
	})
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
