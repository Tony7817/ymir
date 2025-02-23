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
	uId, err := id.GetDecodedUserId(l.ctx)
	if err != nil {
		return nil, err
	}

	pId, err := id.DecodeId(req.ProductId)
	if err != nil {
		return nil, err
	}
	cId, err := id.DecodeId(req.ColorId)
	if err != nil {
		return nil, err
	}
	respb, err := l.svcCtx.ProductRPC.DecreaseProductAmountInCart(l.ctx, &product.DecreaseProductAmountInCartRequest{
		ProductId: pId,
		UserId:    uId,
		ColorId:   cId,
		Size:      req.Size,
	})
	if err != nil {
		return nil, err
	}

	return &types.DecreaseProductAmountInCartResponse{
		ProductCartId: id.EncodeId(respb.ProductCartId),
		Amount:        respb.Amount,
		TotalPrice:    respb.TotalPrice,
	}, nil
}
