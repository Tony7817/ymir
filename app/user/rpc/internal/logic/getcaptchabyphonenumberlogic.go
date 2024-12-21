package logic

import (
	"context"

	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/user"
	"ymir.com/pkg/xerr"

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
	captcha, err := l.svcCtx.CaptchaModel.FindCaptchaByPhonenumber(in.Phonenumber)
	if err != nil {
		return nil, err
	}

	if captcha == nil {
		return nil, xerr.NewErrCode(xerr.WrongCaptchaError)
	}

	return &user.GetCaptchaResponse{
		Captcha:   captcha.VerifyCode,
		CreatedAt: captcha.CreatedAt.Unix(),
	}, nil
}
