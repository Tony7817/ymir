package user

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/user/rpc/user"
	"ymir.com/pkg/id"
	"ymir.com/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDetailLogic {
	return &UserDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDetailLogic) UserDetail(req *types.UserDetailRequest) (*types.UserDetailResponse, error) {
	uId, err := id.GetDecodedUserId(l.ctx)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.UnauthorizedError)
	}

	respb, err := l.svcCtx.UserRPC.GetUserInfo(l.ctx, &user.GetUserInfoRequest{
		Email:       nil,
		Phonenumber: nil,
		UserId:      &uId,
	})
	if err != nil {
		return nil, err
	}

	return &types.UserDetailResponse{
		Id:          id.EncodeId(respb.User.Id),
		Phonenumber: &respb.User.Phonenumber,
		Email:       &respb.User.Email,
		Username:    respb.User.Username,
	}, nil
}
