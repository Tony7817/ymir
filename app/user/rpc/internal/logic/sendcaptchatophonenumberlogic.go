package logic

import (
	"context"

	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendCaptchaToPhonenumberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendCaptchaToPhonenumberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendCaptchaToPhonenumberLogic {
	return &SendCaptchaToPhonenumberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendCaptchaToPhonenumberLogic) SendCaptchaToPhonenumber(in *user.SendCaptchaToPhonenumberRequest) (*user.SendPhonenumberResponse, error) {
	// todo: add your logic here and delete this line

	return &user.SendPhonenumberResponse{}, nil
}
