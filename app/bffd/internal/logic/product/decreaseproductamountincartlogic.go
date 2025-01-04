package product

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/product/rpc/product"
	"ymir.com/pkg/id"

	"github.com/zeromicro/go-zero/core/logx"
)

type DecreaseProductAmountInCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDecreaseProductAmountInCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DecreaseProductAmountInCartLogic {
	return &DecreaseProductAmountInCartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DecreaseProductAmountInCartLogic) DecreaseProductAmountInCart(req *types.DecreaseProductAmountInCartRequest) (resp *types.DecreaseProductAmountInCartResponse, err error) {
	uIdDecoded, err := id.GetDecodedUserId(l.ctx)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.ProductRPC.DecreaseProductAmountInCart(l.ctx, &product.DecreaseProductAmountInCartRequest{
		ProductId: req.ProductId,
		UserId:    uIdDecoded,
		Color:     req.Color,
	})
	if err != nil {
		return nil, err
	}

	return &types.DecreaseProductAmountInCartResponse{}, nil
}
