package product

import (
	"context"

	"github.com/pkg/errors"
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
		idEncoded, err := l.svcCtx.Hash.EncodedId(respb.Products[i].Id)
		if err != nil {
			return nil, errors.Wrapf(err, "[ProductList] failed to encode product id : %d", respb.Products[i].Id)
		}
		ps = append(ps, types.ProductListItem{
			Id:          idEncoded,
			CoverUrl:    respb.Products[i].Coverurl.Url,
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
