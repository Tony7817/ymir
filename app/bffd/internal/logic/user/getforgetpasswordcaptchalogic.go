package user

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetForgetPasswordCaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetForgetPasswordCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetForgetPasswordCaptchaLogic {
	return &GetForgetPasswordCaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetForgetPasswordCaptchaLogic) GetForgetPasswordCaptcha(req *types.GetForgetPasswordCaptchaRequest) (resp *types.GetForgetPasswordCaptchaResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
