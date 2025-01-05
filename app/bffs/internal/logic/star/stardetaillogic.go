package star

import (
	"context"

	"ymir.com/app/bffs/internal/svc"
	"ymir.com/app/bffs/internal/types"
	"ymir.com/app/star/admin/starclient"
	"ymir.com/pkg/id"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type StarDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStarDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StarDetailLogic {
	return &StarDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StarDetailLogic) StarDetail(req *types.StarDetailRequest) (resp *types.StarDetailResponse, err error) {
	sId, err := id.DecodeId(req.Id)
	if err != nil {
		return nil, err
	}
	respb, err := l.svcCtx.StarAdminRPC.StarDetail(l.ctx, &starclient.StarDetailRequest{
		Id: sId,
	})
	if err != nil {
		return nil, err
	}

	var res types.StarDetailResponse
	if err := copier.Copy(&res, &respb); err != nil {
		return nil, err
	}
	res.Id = id.EncodeId(respb.Id)

	return &res, nil
}
