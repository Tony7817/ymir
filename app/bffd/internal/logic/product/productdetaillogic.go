package product

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/product/rpc/product"
	"ymir.com/app/star/rpc/star"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProductDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductDetailLogic {
	return &ProductDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProductDetailLogic) ProductDetail(req *types.ProductDetailRequest) (resp *types.ProductDetailResponse, err error) {
	var pIdDecoded = l.svcCtx.Hash.DecodedId(req.Id)

	p, err := l.svcCtx.ProductRPC.ProductDetail(l.ctx, &product.ProductDetailReqeust{
		Id: pIdDecoded,
	})
	if err != nil {
		return nil, err
	}

	s, err := l.svcCtx.StarRPC.StarDetail(l.ctx, &star.StarDetailRequest{
		Id: p.StarId,
	})
	if err != nil {
		return nil, err
	}

	pIdEncoded, err := l.svcCtx.Hash.EncodedId(p.Id)
	if err != nil {
		return nil, err
	}
	sIdEncoded, err := l.svcCtx.Hash.EncodedId(s.Id)
	if err != nil {
		return nil, err
	}

	var images []types.ProductImage
	var detailImages []types.ProductImage
	for i := 0; i < len(p.Images); i++ {
		images = append(images, types.ProductImage{
			Url: p.Images[i].Url,
		})
	}
	for i := 0; i < len(p.DetailImages); i++ {
		detailImages = append(detailImages, types.ProductImage{
			Url: p.DetailImages[i].Url,
		})
	}
	var res = &types.ProductDetailResponse{
		Id:           pIdEncoded,
		Description:  p.Description,
		Rate:         p.Rate,
		RateCount:    p.ReteCount,
		Price:        p.Price,
		Unit:         p.Unit,
		Size:         p.Size,
		Color:        p.Color,
		Images:       images,
		DetailImages: detailImages,
		SoldNum:      p.SoldNum,
		Detail:       &p.Description,
		StarAvatar:   s.AvatarUrl,
		StarName:     s.Name,
		StarId:       sIdEncoded,
	}

	return res, nil
}
