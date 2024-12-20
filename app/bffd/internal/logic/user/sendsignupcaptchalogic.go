package user

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/user/rpc/user"

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
	resp, err := l.svcCtx.UserRPC.SendCaptchaToEmail(l.ctx, &user.SendCaptchaToEmailRequest{
		Email: *req.Email,
	})
	if err != nil {
		return nil, err
	}

	return &types.SendSignupCaptchaResponse{
		CreatedAt: resp.CreatedAt,
	}, nil
}
