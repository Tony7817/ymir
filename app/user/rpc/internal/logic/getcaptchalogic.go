package logic

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/user"
	"ymir.com/pkg/vars"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCaptchaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCaptchaLogic {
	return &GetCaptchaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCaptchaLogic) GetCaptcha(in *user.GetCaptchaRequest) (*user.GetCaptchaResponse, error) {
	var cacheKey string
	if in.Email != nil {
		cacheKey = vars.GetEmailCapchaCacheKey(*in.Email)
	} else if in.Phonenumber == nil {
		cacheKey = vars.GetPhonenumberCapchaCacheKey(*in.Phonenumber)
	} else {
		return nil, status.Error(codes.InvalidArgument, "email or phonenumber is required")
	}

	captcha, err := l.svcCtx.Redis.GetCtx(l.ctx, cacheKey)
	if err == nil {
		return &user.GetCaptchaResponse{
			Captcha: captcha,
		}, nil
	}

	

	return &user.GetCaptchaResponse{}, nil
}
