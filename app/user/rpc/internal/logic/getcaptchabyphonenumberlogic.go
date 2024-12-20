package logic

import (
	"context"

	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCaptchaByPhonenumberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCaptchaByPhonenumberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCaptchaByPhonenumberLogic {
	return &GetCaptchaByPhonenumberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCaptchaByPhonenumberLogic) GetCaptchaByPhonenumber(in *user.GetCaptchaByPhonenumberRequest) (*user.GetCaptchaResponse, error) {
	// todo: add your logic here and delete this line

	return &user.GetCaptchaResponse{}, nil
}
