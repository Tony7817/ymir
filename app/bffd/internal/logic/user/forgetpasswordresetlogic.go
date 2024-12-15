package user

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ForgetPasswordResetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewForgetPasswordResetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForgetPasswordResetLogic {
	return &ForgetPasswordResetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ForgetPasswordResetLogic) ForgetPasswordReset(req *types.ForgetPasswordRequest) (resp *types.ForgetpasswordResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
