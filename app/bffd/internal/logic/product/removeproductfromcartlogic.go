package product

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"

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
	ret, err := 

	return
}
