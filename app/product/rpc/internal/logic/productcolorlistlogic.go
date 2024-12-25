package logic

import (
	"context"

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
	for i := 0; i < len(cs); i++ {
		res.Colors = append(res.Colors, &product.ProductColorListItem{
			ColorId:  cs[i].Id,
			CoverUrl: cs[i].CoverUrl,
		})
	}

	return res, nil
}
