package user

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetIpAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetIpAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIpAddressLogic {
	return &GetIpAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetIpAddressLogic) GetIpAddress(req *types.GetIpAddressRequest) (resp *types.GetIpAddressResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
