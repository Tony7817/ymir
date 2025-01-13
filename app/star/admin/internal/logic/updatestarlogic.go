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
	var desc sql.NullString
	if in.Description != nil {
		desc.String = *in.Description
		desc.Valid = true
	}
	var starInfo = model.StarPartial{
		Id:          in.Id,
		Name:        in.Name,
		AvatarUrl:   in.AvatarUrl,
		CoverUrl:    in.CoverUrl,
		Description: &desc,
		PosterUrl:   in.PosterUrl,
	}
	_, err := l.svcCtx.StarModel.UpdatePartial(l.ctx, &starInfo)
	if err != nil {
		return nil, err
	}

	return &star.UpdateStarResponse{
		Id: in.Id,
	}, nil
}
