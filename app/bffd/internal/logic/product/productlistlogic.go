package product

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/product/rpc/product"
)

type ProductListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProductListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductListLogic {
	return &ProductListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProductListLogic) ProductList(req *types.ProductListRequest) (*types.ProductListResponse, error) {
	respb, err := l.svcCtx.ProductRPC.ProductList(l.ctx, &product.ProductListRequest{
		StarId:   nil,
		Keyword:  nil,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	var ps []types.ProductListItem
	for i := 0; i < len(respb.Products); i++ {
		ps = append(ps, types.ProductListItem{
			Id:          respb.Products[i].Id,
			CoverUrl:    respb.Products[i].Coverurl,
			Description: respb.Products[i].Description,
			Price:       respb.Products[i].Price,
			Unit:        respb.Products[i].Unit,
		})
	}

	var res = &types.ProductListResponse{
		Total:    respb.Total,
		Products: ps,
	}

	return res, nil
}
