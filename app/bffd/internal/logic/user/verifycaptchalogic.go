package user

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/user/rpc/user"
	"ymir.com/pkg/util"
	"ymir.com/pkg/vars"
	"ymir.com/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyCaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerifyCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyCaptchaLogic {
	return &VerifyCaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerifyCaptchaLogic) VerifyCaptcha(req *types.VerifyCaptchaRequest) (resp *types.VerifyCaptchaResponse, err error) {
	if req.Email != nil {
		if !util.IsEmailValid(*req.Email) {
			return nil, xerr.NewErrCode(xerr.ErrorReuqestParam)
		}
		if err := l.verifyCaptchaByEmail(*req.Email, req.Captcha); err != nil {
			return nil, err
		}
	} else if req.Phonenumber != nil {
		if !util.IsPhonenumberValid(*req.Phonenumber) {
			return nil, xerr.NewErrCode(xerr.ErrorReuqestParam)

		}
		if err := l.verifyCaptchaByPhoneNumber(*req.Phonenumber, req.Captcha); err != nil {
			return nil, err
		}
	} else {
		return nil, xerr.NewErrCode(xerr.ErrorReuqestParam)
	}

	return &types.VerifyCaptchaResponse{
		OK: true,
	}, nil
}

func (l *VerifyCaptchaLogic) verifyCaptchaByEmail(email string, captcha string) (err error) {
	resp, err := l.svcCtx.UserRPC.GetCaptchaByEmail(context.Background(), &user.GetCaptchaByEmailRequest{
		Email: email,
	})
	if err != nil {
		return err
	}
	if captcha == resp.Captcha {
		if _, err = l.svcCtx.UserRPC.DeleteCaptcha(l.ctx, &user.DeleteCaptchaRequest{
			Id:       resp.Id,
			CacheKey: vars.GetCaptchaEmailCacheKey(email),
		}); err != nil {
			return err
		}
	} else {
		return xerr.NewErrCode(xerr.WrongCaptchaError)
	}

	return nil
}

func (l *VerifyCaptchaLogic) verifyCaptchaByPhoneNumber(phoneNumber string, captcha string) (err error) {
	resp, err := l.svcCtx.UserRPC.GetCaptchaByPhonenumber(l.ctx, &user.GetCaptchaByPhonenumberRequest{
		Phonenumber: phoneNumber,
	})
	if err != nil {
		return err
	}
	if captcha == resp.Captcha {
		if _, err = l.svcCtx.UserRPC.DeleteCaptcha(l.ctx, &user.DeleteCaptchaRequest{
			Id:       resp.Id,
			CacheKey: vars.GetCaptchaPhonenumberCacheKey(phoneNumber),
		}); err != nil {
			return err
		}
	} else {
		return xerr.NewErrCode(xerr.WrongCaptchaError)
	}

	return nil
}
