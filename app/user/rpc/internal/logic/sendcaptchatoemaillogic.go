package logic

import (
	"context"

	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendCaptchaToEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendCaptchaToEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendCaptchaToEmailLogic {
	return &SendCaptchaToEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendCaptchaToEmailLogic) SendCaptchaToEmail(in *user.SendCaptchaToEmailRequest) (*user.SendCaptchaToEmailResponse, error) {
	// todo: add your logic here and delete this line

	return &user.SendCaptchaToEmailResponse{}, nil
}
