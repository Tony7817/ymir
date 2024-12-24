package product

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
