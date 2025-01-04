package product

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/product/rpc/product"
	"ymir.com/pkg/id"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveProductFromCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveProductFromCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveProductFromCartLogic {
	return &RemoveProductFromCartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveProductFromCartLogic) RemoveProductFromCart(req *types.RemoveProductFromCartRequest) (resp *types.RemoveProductFromCartResponse, err error) {
	uIdDecoded, err := id.GetDecodedUserId(l.ctx)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.ProductRPC.RemoveProductFromCart(l.ctx, &product.RemoveProductFromCartRequest{
		ProductId: req.ProductId,
		UserId:    uIdDecoded,
		ColorId:   req.ColorId,
	})
	if err != nil {
		return nil, err
	}

	return &types.RemoveProductFromCartResponse{}, nil
}
