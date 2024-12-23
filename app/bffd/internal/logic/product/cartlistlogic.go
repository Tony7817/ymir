package product

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type CartListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCartListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartListLogic {
	return &CartListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CartListLogic) CartList(req *types.ProductCartListRequest) (resp *types.ProductCartListResponse, err error) {
	if req.PageSize > 10 {
		req.PageSize = 10
	}

	uIddeocded := l.svcCtx.Hash.DecodedId(req.UserId)
	respb, err := l.svcCtx.ProductRPC.ProductsInCartList(l.ctx, &product.ProductsInCartListRequest{
		UserId:   uIddeocded,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	var products []types.ProductCartListItem
	for i := 0; i < len(respb.Products); i++ {
		idEncoded, err := l.svcCtx.Hash.EncodedId(respb.Products[i].ProductId)
		if err != nil {
			return nil, err
		}
		products = append(products, types.ProductCartListItem{
			Id:          idEncoded,
			Description: respb.Products[i].Description,
			Price:       respb.Products[i].Price,
			Unit:        respb.Products[i].Unit,
			CoverUrl:    respb.Products[i].CoverUrl,
			Quantity:    respb.Products[i].Amount,
			Amount:      respb.Products[i].Amount,
			Sizes:       respb.Products[i].Sizes,
		})
	}

	return &types.ProductCartListResponse{
		Total:    respb.Total,
		Products: products,
	}, nil
}
