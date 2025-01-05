package recommend

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/star/rpc/star"
	"ymir.com/pkg/id"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecommendStarListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecommendStarListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecommendStarListLogic {
	return &RecommendStarListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecommendStarListLogic) RecommendStarList(req *types.RecommendStarListRequest) (resp *types.RecommendStarListResponse, err error) {
	starspb, err := l.svcCtx.StarRPC.StarList(l.ctx, &star.StarListRequest{
		Keyword:  nil,
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		return nil, err
	}

	var stars []types.RecommendStarListItem
	for i := 0; i < len(starspb.Stars); i++ {
		stars = append(stars, types.RecommendStarListItem{
			StarId:     id.EncodeId(starspb.Stars[i].Id),
			StarAvatar: starspb.Stars[i].AvatarUrl,
			StarName:   starspb.Stars[i].Name,
		})
	}

	var res = &types.RecommendStarListResponse{
		Recommends: stars,
	}

	return res, nil
}
