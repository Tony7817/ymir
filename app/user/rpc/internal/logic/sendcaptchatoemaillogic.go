package logic

import (
	"context"
	"database/sql"
	"strconv"

	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/model"
	"ymir.com/app/user/rpc/user"
	"ymir.com/pkg/util"
	"ymir.com/pkg/vars"
	"ymir.com/pkg/xerr"

	"github.com/pkg/errors"
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
	sendTimesRaw, _ := l.svcCtx.Redis.Get(vars.GetCaptchaEmailSendTimes(in.Email))
	if sendTimesRaw != "" {
		sendTimes, _ := strconv.ParseInt(sendTimesRaw, 10, 64)
		if sendTimes > vars.CaptchaMaxSendTimesPerDay {
			return nil, xerr.NewErrCode(xerr.CaptchaExpireError)
		}
	}
	
	captcha, err := l.svcCtx.CaptchaModel.FindCaptchaByEmail(in.Email)
	if err != nil {
		return nil, err
	}

	if captcha != nil {
		go l.sendCaptchaEmail(in.Email, captcha.VerifyCode)
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

	ret, err := l.svcCtx.CaptchaModel.InsertCaptchaToDbAndCache(newCaptcha)
	if err != nil {
		return nil, err
	}

	lastInsertedI, err := ret.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "[SendCaptchaToEmail] get last insert id failed")
	}

	lastInsertedCaptcha, err := l.svcCtx.CaptchaModel.FindOne(l.ctx, lastInsertedI)
	if err != nil {
		return nil, errors.Wrapf(err, "[SendCaptchaToEmail] find captcha failed, id: %d", lastInsertedI)
	}

	return &user.SendCaptchaToEmailResponse{
		CreatedAt: lastInsertedCaptcha.CreatedAt.Unix(),
	}, nil
}

func (l *SendCaptchaToEmailLogic) sendCaptchaEmail(email string, captcha string) error {
	err := l.svcCtx.EmailClient.SendNoReplayEmail(email, vars.EmailCaptchaSubJect, vars.GetCaptchaEmailTemplate(captcha))
	if err != nil {
		return err
	}

	return nil
}
