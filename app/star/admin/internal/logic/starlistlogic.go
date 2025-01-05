package logic

import (
	"context"

	"ymir.com/app/star/admin/internal/svc"
	"ymir.com/app/star/admin/star"
	"ymir.com/app/star/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
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
	var (
		sl     []model.Star
		count  int64
		offset = (in.Page - 1) * in.PageSize
	)
	err := mr.Finish(func() error {
		var err error
		count, err = l.svcCtx.StarModel.CountStarTotal(l.ctx)
		if err != nil {
			return err
		}
		return nil
	}, func() error {
		var err error
		sl, err = l.svcCtx.StarModel.FindStarList(l.ctx, offset, in.PageSize)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var stars []*star.StarListItem
	for i := 0; i < len(sl); i++ {
		stars = append(stars, &star.StarListItem{
			Id:          sl[i].Id,
			Name:        sl[i].Name,
			Description: sl[i].Description.String,
			AvatarUrl:   sl[i].AvatarUrl,
			Rate:        sl[i].Rate.Float64,
			RateCount:   sl[i].RateCount,
		})
	}

	return &star.StarListResponse{
		Stars: stars,
		Total: count,
	}, nil
}
