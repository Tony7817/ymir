package logic

import (
	"context"
	"database/sql"

	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/model"
	"ymir.com/app/user/rpc/user"
	"ymir.com/pkg/util"
	"ymir.com/pkg/vars"
	"ymir.com/pkg/xerr"

	"github.com/alibabacloud-go/dm-20151123/v2/client"
	"github.com/alibabacloud-go/tea/tea"
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
	if util.IsEmailValid(in.Email) {
		return nil, xerr.NewErrCode(xerr.ReuqestParamError)
	}

	var cacheKey = vars.GetEmailCapchaCacheKey(in.Email)
	err := l.svcCtx.Redis.SetexCtx(l.ctx, cacheKey, "123456", 300)
	if err != nil {
		logx.Errorf("set email captcha to redis failed, err: %+v", err)
		err = nil
	}

	code, err := util.GenerateCpatcha()
	if err != nil {
		return nil, err
	}

	captcha := model.Captcha{
		VerifyCode: code,
		Email: sql.NullString{
			String: in.Email,
			Valid:  true,
		},
		IsDelete: 0,
	}
	_, err = l.svcCtx.CaptchaModel.Insert(l.ctx, &captcha)
	if err != nil {
		return nil, err
	}

	go func() {
		resp, err := l.svcCtx.EmailClient.SingleSendMail(&client.SingleSendMailRequest{
			AccountName:    tea.String(vars.EmailCaptchaSenderName),
			AddressType:    tea.Int32(1),
			ReplyToAddress: tea.Bool(false),
			ToAddress:      tea.String(in.Email),
			Subject:        tea.String(vars.EmailCaptchaSubJect),
			HtmlBody:       tea.String(vars.GetCaptchaEmailTemplate(code)),
			FromAlias:      tea.String("Miss Lover"),
		})
		if err != nil {
			logx.Errorf("[EmailCaptcha] send captcha to email failed, err: %+v", err)
		}
		if resp.StatusCode != nil && (*resp.StatusCode < 200 || *resp.StatusCode >= 300) {
			logx.Errorf("[EmailCaptcha] send captcha to email failed resp: %+v", *resp)
		}
	}()

	return &user.SendCaptchaToEmailResponse{
		Captcha: "123456",
	}, nil
}
