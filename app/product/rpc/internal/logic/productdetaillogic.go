package logic

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"ymir.com/app/product/rpc/internal/svc"
	"ymir.com/app/product/rpc/product"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProductDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductDetailLogic {
	return &ProductDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductDetailLogic) ProductDetail(in *product.ProductDetailReqeust) (*product.ProductDetailResponse, error) {
	p, err := l.svcCtx.ProductModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	var size = strings.Split(p.Size, ",")
	var images []*product.ProductImage
	var detailImages []*product.ProductImage
	logx.Infof("images: %+v", p.Images)
	logx.Infof("detail images: %+v", p.DetailImages)
	if err := json.Unmarshal([]byte(p.Images), &images); err != nil {
		return nil, errors.Wrapf(err, "unmarshal product images failed, err: %+v", err)
	}
	if p.DetailImages.Valid {
		if err := json.Unmarshal([]byte(p.DetailImages.String), &detailImages); err != nil {
			return nil, errors.Wrapf(err, "unmarshal product model images failed, err: %+v", err)
		}
	}
	sort.Slice(images, func(i, j int) bool {
		return images[i].DisPlayRank < images[j].DisPlayRank
	})
	sort.Slice(detailImages, func(i, j int) bool {
		return detailImages[i].DisPlayRank < detailImages[j].DisPlayRank
	})
	var res = &product.ProductDetailResponse{
		Id:           p.Id,
		Description:  p.Description,
		Rate:         p.Rate,
		ReteCount:    p.RateCount,
		Price:        p.Price,
		Unit:         p.Unit,
		Color:        p.Color,
		SoldNum:      p.SoldNum,
		Detail:       p.Detail.String,
		Size:         size,
		StarId:       p.StarId,
		Images:       images,
		DetailImages: detailImages,
	}

	return res, nil
}
