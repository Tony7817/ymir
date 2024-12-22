package user

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/user/rpc/user"
	"ymir.com/pkg/util"
	"ymir.com/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendForgetPasswordCaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendForgetPasswordCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendForgetPasswordCaptchaLogic {
	return &SendForgetPasswordCaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendForgetPasswordCaptchaLogic) SendForgetPasswordCaptcha(req *types.SendForgetPasswordCaptchaRequest) (resp *types.SendForgetPasswordCaptchaResponse, err error) {
	var createdAt int64
	if req.Email != nil {
		if !util.IsEmailValid(*req.Email) {
			return nil, xerr.NewErrCode(xerr.ReuqestParamError)
		}
		resp, err := l.svcCtx.UserRPC.SendCaptchaToEmail(l.ctx, &user.SendCaptchaToEmailRequest{
			Email: *req.Email,
		})
		if err != nil {
			return nil, err
		}
		createdAt = resp.CreatedAt
	} else if req.Phonenumber != nil {
		if !util.IsPhonenumberValid(*req.Phonenumber) {
			return nil, xerr.NewErrCode(xerr.ReuqestParamError)
		}
		resp, err := l.svcCtx.UserRPC.SendCaptchaToPhonenumber(l.ctx, &user.SendCaptchaToPhonenumberRequest{
			Phonenumber: *req.Phonenumber,
		})
		if err != nil {
			return nil, err
		}
		createdAt = resp.CreatedAt
	} else {
		return nil, xerr.NewErrCode(xerr.ReuqestParamError)
	}

	return &types.SendForgetPasswordCaptchaResponse{
		CreatedAt: createdAt,
	}, nil
}
