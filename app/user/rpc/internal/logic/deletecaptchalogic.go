package logic

import (
	"context"

	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCaptchaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCaptchaLogic {
	return &DeleteCaptchaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCaptchaLogic) DeleteCaptcha(in *user.DeleteCaptchaRequest) (*user.DeleteCaptchaResponse, error) {
	

	return &user.DeleteCaptchaResponse{}, nil
}
