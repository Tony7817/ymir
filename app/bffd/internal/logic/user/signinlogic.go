package user

import (
	"context"
	"time"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/user/rpc/user"
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

	respb, err := l.svcCtx.UserRPC.GetUser(l.ctx, &user.GetUserRequest{
		Email:       req.Email,
		Phonenumber: req.Phonenumber,
		UserId:      nil,
	})
	if err != nil {
		return nil, err
	}

	if respb.User == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DataNoExistError), "user not exist")
	}

	var nowDate = time.Now().Unix()
	uIdEncoded, err := l.svcCtx.Hash.EncodedId(respb.User.Id)
	if err != nil {
		return nil, err
	}
	token, err := util.GetJwtToken(l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire, nowDate, uIdEncoded)
	if err != nil {
		return nil, err
	}

	return &types.SigninResponse{
		UserId:      uIdEncoded,
		Username:    respb.User.Username,
		AccessToken: token,
	}, nil
}
