package logic

import (
	"context"

	"ymir.com/app/star/admin/internal/svc"
	"ymir.com/app/star/admin/star"

	"github.com/zeromicro/go-zero/core/logx"
)

type StarDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStarDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StarDetailLogic {
	return &StarDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StarDetailLogic) StarDetail(in *star.StarDetailRequest) (*star.StarDetailResponse, error) {
	// todo: add your logic here and delete this line

	return &star.StarDetailResponse{}, nil
}
