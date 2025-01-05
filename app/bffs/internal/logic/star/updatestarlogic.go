package star

import (
	"context"

	"ymir.com/app/bffs/internal/svc"
	"ymir.com/app/bffs/internal/types"
	"ymir.com/app/star/admin/star"
	"ymir.com/pkg/id"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateStarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateStarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStarLogic {
	return &UpdateStarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateStarLogic) UpdateStar(req *types.UpdateStarRequest) (*types.UpdateStarResponse, error) {
	sId, err := id.DecodeId(req.StarId)
	if err != nil {
		return nil, err
	}

	respb, err := l.svcCtx.StarAdminRPC.UpdateStar(l.ctx, &star.UpdateStarReqeust{
		Id:          sId,
		Name:        req.Name,
		Description: req.Description,
		CoverUrl:    req.CoverUrl,
		AvatarUrl:   req.AvatarUrl,
		PosterUrl:   req.PosterUrl,
	})
	if err != nil {
		return nil, err
	}

	var res types.UpdateStarResponse
	if err := copier.Copy(&res, &respb); err != nil {
		return nil, err
	}
	res.StarId = id.EncodeId(respb.Id)

	return &res, nil
}
