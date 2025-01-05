package logic

import (
	"context"

	"ymir.com/app/star/model"
	"ymir.com/app/star/rpc/internal/svc"
	"ymir.com/app/star/rpc/star"

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
	if in.PageSize > 20 {
		in.PageSize = 20
	}

	var offset = (in.Page - 1) * in.PageSize
	var stars []model.Star
	var total int64
	err := mr.Finish(func() error {
		var err error
		stars, err = l.svcCtx.StarModel.FindStarList(l.ctx, offset, in.PageSize)
		if err != nil {
			return err
		}
		return nil
	},
		func() error {
			var err error
			total, err = l.svcCtx.StarModel.CountStarTotal(l.ctx)
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	var starItems []*star.StarListResponseItem
	for i := 0; i < len(stars); i++ {
		starItems = append(starItems, &star.StarListResponseItem{
			Id:        stars[i].Id,
			Name:      stars[i].Name,
			CoverUrl:  stars[i].CoverUrl,
			AvatarUrl: stars[i].AvatarUrl,
		})
	}
	var res = &star.StarListResponse{
		Stars: starItems,
		Total: total,
	}

	return res, nil
}
