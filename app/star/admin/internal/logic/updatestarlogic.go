package logic

import (
	"context"
	"database/sql"

	"ymir.com/app/star/admin/internal/svc"
	"ymir.com/app/star/admin/star"
	"ymir.com/app/star/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateStarLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateStarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStarLogic {
	return &UpdateStarLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateStarLogic) UpdateStar(in *star.UpdateStarReqeust) (*star.UpdateStarResponse, error) {
	var starInfo = model.Star{Id: in.Id}
	if in.Name != nil {
		starInfo.Name = *in.Name
	}
	if in.AvatarUrl != nil {
		starInfo.AvatarUrl = *in.AvatarUrl
	}
	if in.CoverUrl != nil {
		starInfo.CoverUrl = *in.CoverUrl
	}
	if in.Description != nil {
		starInfo.Description = sql.NullString{
			String: *in.Description,
			Valid:  true,
		}
	}
	if in.PosterUrl != nil {
		starInfo.PosterUrl = *in.PosterUrl
	}
	if err := l.svcCtx.StarModel.Update(l.ctx, &starInfo); err != nil {
		return nil, err
	}

	return &star.UpdateStarResponse{
		Id: in.Id,
	}, nil
}
