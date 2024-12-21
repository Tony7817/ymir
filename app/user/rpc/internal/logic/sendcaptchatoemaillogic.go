package logic

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/model"
	"ymir.com/app/user/rpc/user"
	"ymir.com/pkg/util"
	"ymir.com/pkg/vars"
	"ymir.com/pkg/xerr"

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
	sendTimesRaw, err := l.svcCtx.Redis.Get(vars.GetCaptchaEmailSendTimes(in.Email))
	if err != nil {
		logx.Errorf("[SendCaptchaToEmail] get captcha send times failed, email: %s, err: %+v", in.Email, err)
	}
	if len(sendTimesRaw) != 0 {
		sendTimes, _ := strconv.ParseInt(sendTimesRaw, 10, 64)
		if sendTimes+1 > vars.CaptchaMaxSendTimesPerDay {
			return nil, xerr.NewErrCode(xerr.MaxCaptchaSendTimeError)
		}
		if _, err := l.svcCtx.Redis.Incr(vars.GetCaptchaEmailSendTimes(in.Email)); err != nil {
			logx.Errorf("[SendCaptchaToEmail] incr captcha send times failed, email: %s, err: %+v", in.Email, err)
		}
	} else {
		if err := l.svcCtx.Redis.Setex(vars.GetCaptchaEmailSendTimes(in.Email), "1", vars.CacheExpireIn1d); err != nil {
			logx.Errorf("[SendCaptchaToEmail] set captcha send times failed, email: %s, err: %+v", in.Email, err)
		}
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
