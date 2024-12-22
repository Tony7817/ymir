package logic

import (
	"context"

	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/user"
	"ymir.com/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCaptchaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCaptchaLogic {
	return &DeleteCaptchaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCaptchaLogic) DeleteCaptcha(in *user.DeleteCaptchaRequest) (*user.DeleteCaptchaResponse, error) {
	ret, err := l.svcCtx.CaptchaModel.SoftDelete(in.Id, in.CacheKey)
	if err != nil {
		return nil, err
	}

	updatesRow, err := ret.RowsAffected()
	if err != nil {
		return nil, err
	}

	if updatesRow == 0 {
		return nil, xerr.NewErrCode(xerr.ServerCommonError)
	}

	return &user.DeleteCaptchaResponse{}, nil
}
