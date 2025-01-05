package star

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/star/rpc/star"
	"ymir.com/pkg/id"

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
	starspb, err := l.svcCtx.StarRPC.StarList(l.ctx, &star.StarListRequest{
		Keyword:  req.Keyword,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	var stars []types.StarListItem
	for i := 0; i < len(starspb.Stars); i++ {
		stars = append(stars, types.StarListItem{
			Id:       id.EncodeId(starspb.Stars[i].Id),
			Name:     starspb.Stars[i].Name,
			CoverUrl: starspb.Stars[i].CoverUrl,
		})
	}

	var res = &types.StarListResponse{
		ToTal: starspb.Total,
		Stars: stars,
	}

	return res, nil
}
