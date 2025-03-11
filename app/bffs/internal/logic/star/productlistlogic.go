package star

import (
	"context"

	"ymir.com/app/bffs/internal/svc"
	"ymir.com/app/bffs/internal/types"
	"ymir.com/app/product/admin/product"
	"ymir.com/pkg/id"

	"github.com/zeromicro/go-zero/core/logx"
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
	var sId = new(int64)
	if req.StarId != nil {
		id, err := id.DecodeId(*req.StarId)
		if err != nil {
			return nil, err
		}
		sId = &id
	}
	respb, err := l.svcCtx.ProductAdminRPC.ProductList(l.ctx, &product.ProductListRequest{
		StarId:   sId,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	var pis = make([]types.ProductListItem, len(respb.Products))
	for i := range respb.Products {
		pis[i] = types.ProductListItem{
			Id: id.EncodeId(respb.Products[i].Id),
			Description: respb.Products[i].Description,
			Name: respb.Products[i].Name,
			Rate:      respb.Products[i].Rate,
			RateCount: respb.Products[i].RateCount,
			SoldNum:   0,
			DefaultColor: types.ProductListColorItem{
				CoverUrl: respb.Products[i].DefaultColor.CoverUrl,
				Price:    respb.Products[i].DefaultColor.Price,
				Unit:     respb.Products[i].DefaultColor.Unit,
			},
		}
	}
	return &types.ProductListResponse{
		Total:    respb.Total,
		Products: pis,
	}, nil
}
