package user

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSignupCaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSignupCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSignupCaptchaLogic {
	return &GetSignupCaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSignupCaptchaLogic) GetSignupCaptcha(req *types.GetSignupCaptchaRequest) (resp *types.GetSignupCaptchaResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
