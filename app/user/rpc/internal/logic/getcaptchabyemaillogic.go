package logic

import (
	"context"

	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/user"
	"ymir.com/pkg/vars"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	var cacheKey = vars.GetEmailCapchaCacheKey(in.Email)
	code, err := l.svcCtx.Redis.Get(cacheKey)
	if err == nil && code != "" {
		return &user.GetCaptchaResponse{
			Captcha: code,
		}, nil
	}

	captcha, err := l.svcCtx.CaptchaModel.FindCaptchaByEmail(in.Email)
	if err != nil {
		return nil, err
	}

	if captcha != nil {
		return &user.GetCaptchaResponse{
			Captcha:   captcha.VerifyCode,
			CreatedAt: captcha.CreatedAt.Unix(),
		}, nil
	} else {
		return nil, status.Error(codes.NotFound, "captcha not found")
	}
}
