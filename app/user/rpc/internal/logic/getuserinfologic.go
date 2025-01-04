package logic

import (
	"context"
	"database/sql"
	"errors"

	"ymir.com/app/user/model"
	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/user"
	"ymir.com/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoRequest) (*user.GetUserInfoResponse, error) {
	var usr *model.User
	var err error
	if in.Email != nil {
		usr, err = l.svcCtx.UserModel.FindOneByEmail(l.ctx, sql.NullString{
			String: *in.Email,
			Valid:  true,
		})
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	} else if in.Phonenumber != nil {
		usr, err = l.svcCtx.UserModel.FindOneByPhoneNumber(l.ctx, sql.NullString{
			String: *in.Phonenumber,
			Valid:  true,
		})
	} else if in.UserId != nil {
		usr, err = l.svcCtx.UserModel.FindOne(l.ctx, *in.UserId)
	} else {
		return nil, xerr.NewErrCode(xerr.ReuqestParamError)
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if errors.Is(err, sql.ErrNoRows) {
		return &user.GetUserInfoResponse{
			User: nil,
		}, nil
	}

	return &user.GetUserInfoResponse{
		User: &user.UserInfo{
			Id:          usr.Id,
			Username:    usr.Username,
			AvatarUrl:   usr.AvatarUrl.String,
			Email:       usr.Email.String,
			Phonenumber: usr.PhoneNumber.String,
		},
	}, nil
}
