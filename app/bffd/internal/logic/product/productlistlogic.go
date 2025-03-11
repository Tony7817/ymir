package product

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/product/rpc/product"
	"ymir.com/pkg/id"
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
	for i := range respb.Products {
		ps = append(ps, types.ProductListItem{
			Id:          id.EncodeId(respb.Products[i].Id),
			Description: respb.Products[i].Description,
			DefaultColor: types.ProductListColorItem{
				CoverUrl: respb.Products[i].DefaultColor.CoverUrl,
				Price:    respb.Products[i].DefaultColor.Price,
				Unit:     respb.Products[i].DefaultColor.Unit,
			},
		})
	}

	var res = &types.ProductListResponse{
		Total:    respb.Total,
		Products: ps,
	}

	return res, nil
}
