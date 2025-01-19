package logic

import (
	"context"
	"database/sql"

	"ymir.com/app/star/admin/internal/svc"
	"ymir.com/app/star/admin/star"
	"ymir.com/app/star/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateStarLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateStarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateStarLogic {
	return &CreateStarLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateStarLogic) CreateStar(in *star.CreateStarRequest) (*star.CreateStarResponse, error) {
	var desc sql.NullString
	if in.Description != nil {
		desc.String = *in.Description
		desc.Valid = true
	}
	var starInfo = model.Star{
		Id:   in.StarId,
		Name: in.Name,
		Rate: sql.NullFloat64{
			Float64: 5,
			Valid:   true,
		},
		RateCount:   1,
		AvatarUrl:   in.AvatarUrl,
		Description: desc,
		CoverUrl:    in.CoverUrl,
		PosterUrl:   in.PosterUrl,
	}

	id, err := l.svcCtx.StarModel.InsertStar(l.ctx, &starInfo)
	if err != nil {
		return nil, err
	}

	return &star.CreateStarResponse{
		Id: id,
	}, nil
}
