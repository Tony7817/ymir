package user

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/user/rpc/user"
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
		if err := l.verifyCaptchaByEmail(*req.Email, req.Captcha); err != nil {
			return nil, err
		}
	} else if req.Phonenumber != nil {
		if err := l.verifyCaptchaByPhoneNumber(*req.Phonenumber, req.Captcha); err != nil {
			return nil, err
		}
	} else {
		return nil, xerr.NewErrCode(xerr.ReuqestParamError)
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
	logx.Debugf("resp captcha: %s", resp.Captcha)
	logx.Debugf("captcha: %s", captcha)
	if captcha == resp.Captcha {
		if _, err = l.svcCtx.UserRPC.DeleteCaptcha(l.ctx, &user.DeleteCaptchaRequest{
			Email: &email,
			Id:    resp.Id,
		}); err != nil {
			return err
		}
		return nil
	} else {
		return xerr.NewErrCode(xerr.WrongCaptchaError)
	}
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
			Phonenumber: &phoneNumber,
			Id:          resp.Id,
		}); err != nil {
			return err
		}
		return nil
	} else {
		return xerr.NewErrCode(xerr.WrongCaptchaError)
	}
}
