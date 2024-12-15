package user

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyCaptchaWithAccountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerifyCaptchaWithAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyCaptchaWithAccountLogic {
	return &VerifyCaptchaWithAccountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerifyCaptchaWithAccountLogic) VerifyCaptchaWithAccount(req *types.VerifyCaptchaRequest) (resp *types.VerifyCaptchaResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
