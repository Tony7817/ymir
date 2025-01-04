package logic

import (
	"context"

	"ymir.com/app/star/admin/internal/svc"
	"ymir.com/app/star/admin/star"

	"github.com/zeromicro/go-zero/core/logx"
)

type StarListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStarListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StarListLogic {
	return &StarListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StarListLogic) StarList(in *star.StarListRequest) (*star.StarListResponse, error) {
	// todo: add your logic here and delete this line

	return &star.StarListResponse{}, nil
}
