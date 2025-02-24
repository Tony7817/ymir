package user

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/user/rpc/user"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CheckIfUserSignedUpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckIfUserSignedUpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckIfUserSignedUpLogic {
	return &CheckIfUserSignedUpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckIfUserSignedUpLogic) CheckIfUserSignedUp(req *types.CheckIfUserSignedUpRequest) (*types.CheckIfUserSignedUpResponse, error) {
	var respb *user.GetUserInfoResponse
	var err error
	if req.Email != nil && *req.Email != "" {
		respb, err = l.svcCtx.UserRPC.GetUserInfo(l.ctx, &user.GetUserInfoRequest{
			Email: req.Email,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "get user info by email failed, email: %s", *req.Email)
		}
	} else if req.PhoneNumber != nil && *req.PhoneNumber != "" {
		respb, err = l.svcCtx.UserRPC.GetUserInfo(l.ctx, &user.GetUserInfoRequest{
			Phonenumber: req.PhoneNumber,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "get user info by phone number failed, phone number: %s", *req.PhoneNumber)
		}
	} else {
		return nil, errors.New("email or phone number must be set")
	}

	if respb.User == nil {
		return &types.CheckIfUserSignedUpResponse{
			OK: false,
		}, nil
	}

	return &types.CheckIfUserSignedUpResponse{
		OK: true,
	}, nil
}
