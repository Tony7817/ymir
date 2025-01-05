package logic

import (
	"context"

	"ymir.com/app/star/rpc/internal/svc"
	"ymir.com/app/star/rpc/star"

	"github.com/pkg/errors"
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
	l.Logger.Infof("star id:%v", in.Id)
	s, err := l.svcCtx.StarModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, errors.Wrapf(err, "[StarDetail] find star by id failed, err: %+v", err)
	}

	var res = &star.StarDetailResponse{
		Id:          s.Id,
		Name:        s.Name,
		Description: s.Description.String,
		CoverUrl:    s.CoverUrl,
		AvatarUrl:   s.AvatarUrl,
		PosterUrl:   s.PosterUrl,
		Rate:        s.Rate.Float64,
	}

	return res, nil
}
