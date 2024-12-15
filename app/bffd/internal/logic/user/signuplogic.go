package user

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/user/rpc/user"
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
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ReuqestParamError), "email or phonenumber must be set")
	}
	// If user exist
	respb, err := l.svcCtx.UserRPC.GetUser(l.ctx, &user.GetUserRequest{
		Email:       req.Email,
		Phonenumber: req.Phonenumber,
		UserId:      nil,
	})
	if err != nil {
		return nil, err
	}

	if respb.User != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.UserAlreadyExistError), "user exist already")
	}

	username, err := util.Mask(*req.Email)
	if err != nil {
		return nil, err
	}
	passwordHash, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	var userId int64
	// sign up
	if req.Email != nil {
		respb, err := l.svcCtx.UserRPC.WriteUserInDBWithEmail(l.ctx, &user.WriteUserInDBWithEmailRequest{
			Email:        *req.Email,
			PasswordHash: passwordHash,
			UserName:     username,
		})
		if err != nil {
			return nil, err
		}
		userId = respb.UserId
	} else if req.Phonenumber != nil {
		respb, err := l.svcCtx.UserRPC.WriteUserInDBWithPhonenumber(l.ctx, &user.WriteUserInDBWithPhonenumberRequest{
			Phonenumber:  *req.Phonenumber,
			PasswordHash: passwordHash,
			UserName:     username,
		})
		if err != nil {
			return nil, err
		}
		userId = respb.UserId
	} else {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ReuqestParamError), "email or phonenumber must be set")
	}

	userIdEncoded, err := l.svcCtx.Hash.EncodedId(userId)
	if err != nil {
		return nil, err
	}

	return &types.SignupResponse{
		UserId: userIdEncoded,
	}, nil
}
