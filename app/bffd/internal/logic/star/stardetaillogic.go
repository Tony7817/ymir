package star

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/product/rpc/product"
	"ymir.com/app/star/rpc/star"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type StarDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStarDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StarDetailLogic {
	return &StarDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StarDetailLogic) StarDetail(req *types.StarDetailRequest) (*types.StarDetailResponse, error) {
	var (
		starId     *int64
		products   *product.ProductListResponse
		starDetail *star.StarDetailResponse
	)

	var starIdDecoded = l.svcCtx.Hash.DecodedId(req.Id)
	starId = &starIdDecoded

	err := mr.Finish(
		func() error {
			var err error
			if products, err = l.svcCtx.ProductRPC.ProductList(l.ctx, &product.ProductListRequest{
				StarId:   starId,
				Keyword:  nil,
				Page:     1,
				PageSize: 20,
			}); err != nil {
				return err
			}
			return nil
		},
		func() error {
			var err error
			if starDetail, err = l.svcCtx.StarRPC.StarDetail(l.ctx, &star.StarDetailRequest{
				Id: *starId,
			}); err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	starIdEncoded, err := l.svcCtx.Hash.EncodedId(starDetail.Id)
	if err != nil {
		return nil, err
	}

	var ps []types.ProductListItem
	for i := 0; i < len(products.Products); i++ {
		idEncoded, err := l.svcCtx.Hash.EncodedId(products.Products[i].Id)
		if err != nil {
			return nil, err
		}
		ps = append(ps, types.ProductListItem{
			Id:          idEncoded,
			CoverUrl:    products.Products[i].Coverurl.Url,
			Description: products.Products[i].Description,
			Price:       products.Products[i].Price,
			Unit:        products.Products[i].Unit,
		})
	}

	var res = &types.StarDetailResponse{
		Id:          starIdEncoded,
		Name:        starDetail.Name,
		Description: starDetail.Description,
		CoverUrl:    starDetail.CoverUrl,
		AvatarUrl:   starDetail.AvatarUrl,
		PosterUrl:   starDetail.PosterUrl,
		Products: types.ProductListResponse{
			Total:    int64(len(ps)),
			Products: ps,
		},
	}

	return res, nil
}
