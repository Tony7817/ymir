package product

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/product/rpc/product"
	"ymir.com/pkg/util"
	"ymir.com/pkg/vars"
	"ymir.com/pkg/xerr"

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
	userIdRaw, ok := l.ctx.Value(vars.UserIdKey).(string)
	if !ok {
		return nil, xerr.NewErrCode(xerr.UnauthorizedError)
	}

	var uIdDoceded = l.svcCtx.Hash.DecodedId(userIdRaw)

	if req.PageSize > 10 {
		req.PageSize = 10
	}

	respb, err := l.svcCtx.ProductRPC.ProductsInCartList(l.ctx, &product.ProductsInCartListRequest{
		UserId:   uIdDoceded,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	var (
		products []types.ProductCartListItem
	)
	for i := 0; i < len(respb.Products); i++ {
		pIdEncoded, err := l.svcCtx.Hash.EncodedId(respb.Products[i].ProductId)
		if err != nil {
			return nil, err
		}
		products = append(products, types.ProductCartListItem{
			ProductId:   pIdEncoded,
			Description: respb.Products[i].Description,
			Price:       respb.Products[i].Price,
			Unit:        respb.Products[i].Unit,
			CoverUrl:    respb.Products[i].CoverUrl,
			Amount:      respb.Products[i].Amount,
			Size:        respb.Products[i].Sizes,
			TotalPrice:  util.MutiplyAndRound(float64(respb.Products[i].Amount), respb.Products[i].Price),
		})
	}

	return &types.ProductCartListResponse{
		Total:    respb.Total,
		Products: products,
	}, nil
}
