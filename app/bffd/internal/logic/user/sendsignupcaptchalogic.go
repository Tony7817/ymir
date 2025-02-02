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

type SendSignupCaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendSignupCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendSignupCaptchaLogic {
	return &SendSignupCaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendSignupCaptchaLogic) SendSignupCaptcha(req *types.SendSignupCaptchaRequest) (*types.SendSignupCaptchaResponse, error) {
	var createdAt int64
	if req.Email != nil {
		resp, err := l.svcCtx.UserRPC.SendCaptchaToEmail(l.ctx, &user.SendCaptchaToEmailRequest{
			Email: *req.Email,
		})
		if err != nil {
			return nil, err
		}
		createdAt = resp.CreatedAt
	} else if req.Phonenumber != nil {
		if !util.IsPhonenumberValid(*req.Phonenumber) {
			return nil, xerr.NewErrCode(xerr.ErrorReuqestParam)
		}
		resp, err := l.svcCtx.UserRPC.SendCaptchaToPhonenumber(l.ctx, &user.SendCaptchaToPhonenumberRequest{
			Phonenumber: *req.Phonenumber,
		})
		if err != nil {
			return nil, err
		}
		createdAt = resp.CreatedAt
	} else {
		return nil, xerr.NewErrCode(xerr.ErrorReuqestParam)
	}

	return &types.SendSignupCaptchaResponse{
		CreatedAt: createdAt,
	}, nil
}
