package logic

import (
	"context"
	"database/sql"

	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/model"
	"ymir.com/app/user/rpc/user"
	"ymir.com/pkg/vars"

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
	var cacheKey = vars.GetPhonenumberCapchaCacheKey(in.Phonenumber)

	err := l.svcCtx.Redis.SetexCtx(l.ctx, cacheKey, "123456", vars.CacheExpireIn300s)
	if err != nil {
		logx.Errorf("set phonenumber captcha to redis failed, err: %+v", err)
		err = nil
	}

	var capcha = model.Captcha{
		VerifyCode: "123456",
		PhoneNumber: sql.NullString{
			String: in.Phonenumber,
			Valid:  true,
		},
	}
	_, err = l.svcCtx.CaptchaModel.Insert(l.ctx, &capcha)
	if err != nil {
		return nil, err
	}

	return &user.SendPhonenumberResponse{
		Captcha: "123456",
	}, nil
}
