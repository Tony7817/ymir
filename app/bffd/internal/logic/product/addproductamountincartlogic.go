package product

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/product/rpc/product"
	"ymir.com/pkg/id"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddProductAmountInCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddProductAmountInCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddProductAmountInCartLogic {
	return &AddProductAmountInCartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddProductAmountInCartLogic) AddProductAmountInCart(req *types.AddProductAmountInCartRequest) (resp *types.AddProductAmountInCartResponse, err error) {
	uIdDecoded, err := id.GetDecodedUserId(l.ctx)
	if err != nil {
		return nil, err
	}

	pId, err := id.DecodeId(req.ProductId)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.ProductRPC.AddProductAmountInCart(l.ctx, &product.AddProductAmountInCartRequest{
		ProductId: pId,
		UserId:    uIdDecoded,
	})
	if err != nil {
		return nil, err
	}

	return &types.AddProductAmountInCartResponse{}, nil
}
