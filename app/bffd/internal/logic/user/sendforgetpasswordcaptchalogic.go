package user

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendForgetPasswordCaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendForgetPasswordCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendForgetPasswordCaptchaLogic {
	return &SendForgetPasswordCaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendForgetPasswordCaptchaLogic) SendForgetPasswordCaptcha(req *types.SendForgetPasswordCaptchaRequest) (resp *types.SendForgetPasswordCaptchaResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
