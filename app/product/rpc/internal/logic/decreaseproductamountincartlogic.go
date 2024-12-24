package logic

import (
	"context"

	"ymir.com/app/product/rpc/internal/svc"
	"ymir.com/app/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type DecreaseProductAmountInCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDecreaseProductAmountInCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DecreaseProductAmountInCartLogic {
	return &DecreaseProductAmountInCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DecreaseProductAmountInCartLogic) DecreaseProductAmountInCart(in *product.DecreaseProductAmountInCartRequest) (*product.AddProductAmountInCartResponse, error) {
	if err := l.svcCtx.ProductCartModel.DescProductAmount(l.ctx, in.UserId, in.ProductId, in.Color); err != nil {
		return nil, err
	}

	return &product.AddProductAmountInCartResponse{}, nil
}
