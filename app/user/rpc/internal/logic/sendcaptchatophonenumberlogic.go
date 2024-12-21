package logic

import (
	"context"
	"database/sql"
	"time"

	"ymir.com/app/user/rpc/internal/common"
	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/model"
	"ymir.com/app/user/rpc/user"
	"ymir.com/pkg/util"
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

func (l *SendCaptchaToPhonenumberLogic) SendCaptchaToPhonenumber(in *user.SendCaptchaToPhonenumberRequest) (*user.SendCaptchaToPhonenumberResponse, error) {
	if err := common.CheckCapchaSendTimes(l.svcCtx.Redis, vars.GetCaptchaPhonenumberSendTimesKey(in.Phonenumber)); err != nil {
		return nil, err
	}
	var resp = &user.SendCaptchaToPhonenumberResponse{
		CreatedAt: time.Now().UTC().Unix(),
	}

	captcha, err := l.svcCtx.CaptchaModel.FindCaptchaByPhonenumber(in.Phonenumber)
	if err != nil {
		return nil, err
	}

	if captcha != nil {
		go util.SendCaptchaToPhonenumber(in.Phonenumber, captcha.VerifyCode)
		return resp, nil
	}

	code, err := util.GenerateCpatcha()
	if err != nil {
		code = vars.CaptchaCodeAnyWay
	}

	var newCaptcha = model.Captcha{
		VerifyCode: code,
		PhoneNumber: sql.NullString{
			String: in.Phonenumber,
			Valid:  true,
		},
	}

	_, err = l.svcCtx.CaptchaModel.InsertCaptchaToDbAndCache(newCaptcha)
	if err != nil {
		return nil, err
	}

	go util.SendCaptchaToPhonenumber(in.Phonenumber, code)

	return resp, nil
}
