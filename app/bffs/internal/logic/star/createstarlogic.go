package star

import (
	"context"

	"ymir.com/app/bffs/internal/svc"
	"ymir.com/app/bffs/internal/types"
	"ymir.com/app/star/admin/star"
	"ymir.com/pkg/id"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateStarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateStarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateStarLogic {
	return &CreateStarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateStarLogic) CreateStar(req *types.CreteStarRequest) (resp *types.CreateStarResponse, err error) {
	respb, err := l.svcCtx.StarAdminRPC.CreateStar(l.ctx, &star.CreateStarRequest{
		Name:        req.Name,
		Description: req.Description,
		CoverUrl:    req.CoverUrl,
		AvatarUrl:   req.AvatarUrl,
		PosterUrl:   req.PosterUrl,
	})
	if err != nil {
		return nil, err
	}

	return &types.CreateStarResponse{
		StarId: id.EncodeId(respb.Id),
	}, nil
}
