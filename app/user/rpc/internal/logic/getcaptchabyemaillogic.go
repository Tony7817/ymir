package logic

import (
	"context"

	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/user"
	"ymir.com/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetCaptchaByEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCaptchaByEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCaptchaByEmailLogic {
	return &GetCaptchaByEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCaptchaByEmailLogic) GetCaptchaByEmail(in *user.GetCaptchaByEmailRequest) (*user.GetCaptchaResponse, error) {
	captcha, err := l.svcCtx.CaptchaModel.FindCaptchaByEmail(in.Email)
	if err != nil {
		return nil, err
	}

	if captcha == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.WrongCaptchaError), "[GetCaptchaByEmail] captcha not found by email:%s", in.Email)
	}

	return &user.GetCaptchaResponse{
		Id:        captcha.Id,
		Captcha:   captcha.VerifyCode,
		CreatedAt: captcha.CreatedAt.Unix(),
	}, nil
}
