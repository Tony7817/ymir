package user

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/user/rpc/user"
	"ymir.com/pkg/id"
	"ymir.com/pkg/util"
	"ymir.com/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type SignupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignupLogic {
	return &SignupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignupLogic) Signup(req *types.SignupRequest) (*types.SignupResponse, error) {
	if req.Email == nil && req.Phonenumber == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ErrorReuqestParam), "email or phonenumber must be set")
	}
	// If user exist
	respb, err := l.svcCtx.UserRPC.GetUserInfo(l.ctx, &user.GetUserInfoRequest{
		Email:       req.Email,
		Phonenumber: req.Phonenumber,
	})
	if err != nil {
		return nil, err
	}

	if respb.User != nil {
		return nil, xerr.NewErrCode(xerr.UserAlreadyExistError)
	}

	var username string
	if req.Email != nil {
		username, err = util.MaskEmail(*req.Email)
	} else {
		username, err = util.MaskPhonenumber(*req.Phonenumber)
	}
	if err != nil {
		return nil, err
	}
	passwordHash, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// sign up
	respbWriteUser, err := l.svcCtx.UserRPC.WriteUserLocalInDB(l.ctx, &user.WriteUserLocalRequest{
		Email:        req.Email,
		Phonenumber:  req.Phonenumber,
		PasswordHash: passwordHash,
		UserName:     username,
	})
	if err != nil {
		return nil, err
	}

	return &types.SignupResponse{
		UserId: id.EncodeId(respbWriteUser.UserId),
	}, nil
}
