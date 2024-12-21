package logic

import (
	"context"
	"database/sql"
	"errors"

	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/model"
	"ymir.com/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *user.GetUserRequest) (*user.GetUserResponse, error) {
	var usr *model.User
	var err error
	if in.UserId != nil {
		usr, err = l.svcCtx.UserModel.FindOne(l.ctx, *in.UserId)
	} else if in.Email != nil {
		usr, err = l.svcCtx.UserModel.FindOneByEmail(l.ctx, sql.NullString{
			String: *in.Email,
			Valid:  true,
		})
	} else if in.Phonenumber != nil {
		usr, err = l.svcCtx.UserModel.FindOneByPhoneNumber(l.ctx, sql.NullString{
			String: *in.Phonenumber,
			Valid:  true,
		})
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return &user.GetUserResponse{
			User: nil,
		}, nil
	}

	var res = &user.GetUserResponse{
		User: &user.UserInfo{
			Username:     usr.Username,
			AvatarUrl:    usr.AvatarUrl.String,
			PasswordHash: usr.PasswordHash,
			IsActivated:  usr.IsActivated,
		},
	}

	return res, nil
}
