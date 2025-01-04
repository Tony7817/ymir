package logic

import (
	"context"
	"database/sql"
	"time"

	"ymir.com/app/user/model"
	"ymir.com/app/user/rpc/internal/common"
	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/user"
	"ymir.com/pkg/id"
	"ymir.com/pkg/util"
	"ymir.com/pkg/vars"

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
	if err := common.CheckCapchaSendTimes(l.svcCtx.Redis, vars.GetCaptchaEmailSendTimesKey(in.Email)); err != nil {
		return nil, err
	}
	var resp = &user.SendCaptchaToEmailResponse{
		CreatedAt: time.Now().UTC().Unix(),
	}

	// find captcha in redis and mysql
	captcha, err := l.svcCtx.CaptchaModel.FindCaptchaByEmail(in.Email)
	if err != nil {
		return nil, err
	}

	if captcha != nil {
		go l.sendCaptchaEmail(in.Email, captcha.VerifyCode)
		return resp, nil
	}

	code, err := util.GenerateCpatcha()
	if err != nil {
		code = vars.CaptchaCodeAnyWay
	}

	var newCaptcha = model.Captcha{
		Id:         id.SF.GenerateID(),
		VerifyCode: code,
		Email: sql.NullString{
			String: in.Email,
			Valid:  true,
		},
	}

	_, err = l.svcCtx.CaptchaModel.InsertCaptchaToDbAndCache(newCaptcha)
	if err != nil {
		return nil, err
	}

	go l.sendCaptchaEmail(in.Email, code)

	return resp, nil
}

func (l *SendCaptchaToEmailLogic) sendCaptchaEmail(email string, captcha string) error {
	err := l.svcCtx.EmailClient.SendNoReplyEmail(email, vars.EmailCaptchaSubJect, vars.GetCaptchaEmailTemplate(captcha))
	if err != nil {
		return err
	}

	return nil
}
