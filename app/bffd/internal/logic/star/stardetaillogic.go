package star

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/product/rpc/product"
	"ymir.com/app/star/rpc/star"
	"ymir.com/pkg/id"

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
		products   *product.ProductListResponse
		starDetail *star.StarDetailResponse
		ps         []types.ProductListItem
	)
	sId, err := id.DecodeId(req.Id)
	if err != nil {
		return nil, err
	}

	err = mr.Finish(
		func() error {
			var err error
			if products, err = l.svcCtx.ProductRPC.ProductList(l.ctx, &product.ProductListRequest{
				StarId:   &sId,
				Keyword:  nil,
				Page:     1,
				PageSize: 20,
			}); err != nil {
				return err
			}
			for i := 0; i < len(products.Products); i++ {
				ps = append(ps, types.ProductListItem{
					Id:          id.EncodeId(products.Products[i].Id),
					Description: products.Products[i].Description,
					DefaultColor: types.ProductListColorItem{
						CoverUrl: products.Products[i].DefaultColor.CoverUrl,
						Price:    products.Products[i].DefaultColor.Price,
						Unit:     products.Products[i].DefaultColor.Unit,
					},
				})
			}

			return nil
		},
		func() error {
			var err error
			if starDetail, err = l.svcCtx.StarRPC.StarDetail(l.ctx, &star.StarDetailRequest{
				Id: sId,
			}); err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	var res = &types.StarDetailResponse{
		Id:          id.EncodeId(starDetail.Id),
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
