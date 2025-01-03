package star

import (
	"context"

	"ymir.com/app/bffs/internal/svc"
	"ymir.com/app/bffs/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StarListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStarListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StarListLogic {
	return &StarListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StarListLogic) StarList(req *types.StarListRequest) (resp *types.StarListResponse, err error) {

	return nil, nil
}
