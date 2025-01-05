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
	s, err := l.svcCtx.StarModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &star.StarDetailResponse{
		Id:          s.Id,
		Name:        s.Name,
		Description: s.Description.String,
		CoverUrl:    s.CoverUrl,
		AvatarUrl:   s.AvatarUrl,
		PosterUrl:   s.PosterUrl,
		Rate:        s.Rate.Float64,
		RateCount:   s.RateCount,
	}, nil
}
