package logic

import (
	"context"
	"strings"

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
	if err != nil {
		return nil, err
	}

	var (
		products []*product.ProductsInCartListItem
		total    int64
	)
	mr.Finish(
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
		p, err := l.svcCtx.ProductModel.FindOne(l.ctx, pcart.ProductId)
		if err != nil {
			cancel(err)
			return
		}
		writer.Write(&product.ProductsInCartListItem{
			ProductId:   p.Id,
			Amount:      pcart.Amount,
			Sizes:       strings.Split(pcart.Size, ","),
			Description: p.Description,
			Price:       p.Price,
			Unit:        p.Unit,
			CoverUrl:    p.CoverUrl,
			QUantity:    pcart.Amount,
		})
	}, func(pipe <-chan *product.ProductsInCartListItem, writer mr.Writer[[]*product.ProductsInCartListItem], cancel func(error)) {
		pitems := []*product.ProductsInCartListItem{}
		for item := range pipe {
			pitems = append(pitems, item)
		}
		writer.Write(pitems)
	})
}
