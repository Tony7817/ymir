package product

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/product/rpc/product"
	"ymir.com/pkg/id"

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
	uId, err := id.GetDecodedUserId(l.ctx)
	if err != nil {
		return nil, err
	}

	if req.PageSize > 10 {
		req.PageSize = 10
	}

	respb, err := l.svcCtx.ProductRPC.ProductsInCartList(l.ctx, &product.ProductsInCartListRequest{
		UserId:   uId,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	var (
		products = make([]types.ProductCartListItem, 0)
	)
	for i := range respb.Products {
		products = append(products, types.ProductCartListItem{
			ProductCartId: id.EncodeId(respb.Products[i].ProductCartId),
			ProductId:     id.EncodeId(respb.Products[i].ProductId),
			StarId:        id.EncodeId(respb.Products[i].StarId),
			ColorId:       id.EncodeId(respb.Products[i].ColorId),
			Description:   respb.Products[i].Description,
			Price:         respb.Products[i].Price,
			Unit:          respb.Products[i].Unit,
			CoverUrl:      respb.Products[i].CoverUrl,
			Amount:        respb.Products[i].Amount,
			Size:          respb.Products[i].Sizes,
			TotalPrice:    respb.Products[i].Amount * respb.Products[i].Price,
			Stock:         respb.Products[i].Stock,
			Color:         respb.Products[i].Color,
		})
	}

	return &types.ProductCartListResponse{
		Total:    respb.Total,
		Products: products,
	}, nil
}
