package product

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddProductToCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddProductToCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddProductToCartLogic {
	return &AddProductToCartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddProductToCartLogic) AddProductToCart(req *types.AddProductToCartRequest) (resp *types.AddProductToCartResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
