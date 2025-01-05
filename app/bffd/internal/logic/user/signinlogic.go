package user

import (
	"context"
	"time"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/user/model"
	"ymir.com/app/user/rpc/user"
	"ymir.com/pkg/id"
	"ymir.com/pkg/util"
	"ymir.com/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type SigninLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSigninLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SigninLogic {
	return &SigninLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SigninLogic) Signin(req *types.SigninRequest) (resp *types.SigninResponse, err error) {
	if req.Email != nil {
		if !util.IsEmailValid(*req.Email) {
			return nil, xerr.NewErrCode(xerr.ReuqestParamError)
		}
	} else if req.Phonenumber != nil {
		if !util.IsPhonenumberValid(*req.Phonenumber) {
			return nil, xerr.NewErrCode(xerr.ReuqestParamError)
		}
	} else {
		return nil, xerr.NewErrCode(xerr.ReuqestParamError)
	}

	respb, err := l.svcCtx.UserRPC.GetUserInfo(l.ctx, &user.GetUserInfoRequest{
		Email:       req.Email,
		Phonenumber: req.Phonenumber,
		UserId:      nil,
	})
	if err != nil {
		return nil, err
	}

	if respb.User == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.UserNotExistedError), "user not exist")
	}

	if respb.User.Type == model.UserTypeGoogle {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.ErrorWrongPassword), "google user cannot use password")
	}

	respbUserLocal, err := l.svcCtx.UserRPC.GetUserLocal(l.ctx, &user.GetUserLocalRequest{
		UserId: respb.User.Id,
	})
	if err != nil {
		return nil, err
	}

	if !util.CheckPassword(req.Password, respbUserLocal.User.PasswordHash) {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.ErrorWrongPassword), "wrong password")
	}

	var nowDate = time.Now().Unix()
	token, err := util.GetJwtToken(l.svcCtx.Config.Auth.AccessSecret, nowDate, l.svcCtx.Config.Auth.AccessExpire, respb.User.Id)
	if err != nil {
		return nil, err
	}

	return &types.SigninResponse{
		UserId:      id.EncodeId(respb.User.Id),
		Username:    respb.User.Username,
		AccessToken: token,
		AvatarUrl:   respb.User.AvatarUrl,
	}, nil
}
