package user

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/user/rpc/user"
	"ymir.com/pkg/id"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOssSTSTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOssSTSTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOssSTSTokenLogic {
	return &GetOssSTSTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOssSTSTokenLogic) GetOssSTSToken(req *types.OssTokenRequest) (*types.OssTokenResponse, error) {
	uId, err := id.GetDecodedUserId(l.ctx)
	if err != nil {
		return nil, err
	}
	respb, err := l.svcCtx.UserRPC.GetOssStsToken(l.ctx, &user.GetOssStsTokenRequest{
		UserId: uId,
	})
	if err != nil {
		return nil, err
	}

	var res = types.OssTokenResponse{
		AccessKeyId:     respb.AccessKeyId,
		AccessKeySecret: respb.AccessKeySecret,
		Expiration:      respb.Expiration,
		SecurityToken:   respb.SecurityToken,
	}

	return &res, nil
}
