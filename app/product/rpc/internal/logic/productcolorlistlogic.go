package logic

import (
	"context"
	"strings"

	"ymir.com/app/product/rpc/internal/svc"
	"ymir.com/app/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductColorListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductColorListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductColorListLogic {
	return &ProductColorListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductColorListLogic) ProductColorList(in *product.ProductColorListRequest) (*product.ProductColorListResponse, error) {
	cs, err := l.svcCtx.ProductColorModel.FindColorListByPId(l.ctx, in.ProductId)
	if err != nil {
		return nil, err
	}

	var res = &product.ProductColorListResponse{}
	for i := range cs {
		res.Colors = append(res.Colors, &product.ProductColor{
			Id:             cs[i].Id,
			ProductId:      cs[i].ProductId,
			Name:           cs[i].Color,
			Images:         strings.Split(cs[i].Images, ","),
			DetailImages:   strings.Split(cs[i].DetailImages, ","),
			Price:          cs[i].Price,
			CoverUrl:       cs[i].CoverUrl,
			Unit:           cs[i].Unit,
			AvaliableSizes: strings.Split(cs[i].AvailableSize, ","),
			IsDefault:      cs[i].IsDefault == 1,
		})
	}

	return res, nil
}
