package star

import (
	"context"

	"ymir.com/app/bffs/internal/svc"
	"ymir.com/app/bffs/internal/types"
	"ymir.com/app/product/admin/productclient"
	"ymir.com/app/star/admin/star"
	"ymir.com/app/star/admin/starclient"
	"ymir.com/pkg/id"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
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

func (l *StarListLogic) StarList(req *types.StarListRequest) (*types.StarListResponse, error) {
	starpb, err := l.svcCtx.StarAdminRPC.StarList(l.ctx, &starclient.StarListRequest{
		Keyword:  req.Keyword,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	res, err := mr.MapReduce(func(source chan<- *star.StarListItem) {
		for i := 0; i < len(starpb.Stars); i++ {
			source <- starpb.Stars[i]
		}
	}, func(item *star.StarListItem, writer mr.Writer[*types.StarListItem], cancel func(error)) {
		pb, err := l.svcCtx.ProductAdminRPC.ProductCount(l.ctx, &productclient.ProductCountRequest{
			StarId: &item.Id,
		})
		if err != nil {
			cancel(err)
			return
		}
		var star = types.StarListItem{
			Id:           id.EncodeId(item.Id),
			Name:         item.Name,
			AvatarUrl:    item.AvatarUrl,
			Description:  item.Description,
			ProductTotal: pb.Count,
			Rate:         item.Rate,
			RateCount:    item.RateCount,
		}
		writer.Write(&star)
	}, func(pipe <-chan *types.StarListItem, writer mr.Writer[[]types.StarListItem], cancel func(error)) {
		var stars []types.StarListItem
		for item := range pipe {
			stars = append(stars, *item)
		}
		writer.Write(stars)
	})
	if err != nil {
		return nil, err
	}

	var resp = types.StarListResponse{
		ToTal: starpb.Total,
		Stars: res,
	}

	return &resp, nil
}
