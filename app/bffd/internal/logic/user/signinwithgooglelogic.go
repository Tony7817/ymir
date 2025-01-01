package user

import (
	"context"
	"time"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/user/rpc/user"
	"ymir.com/pkg/google"
	"ymir.com/pkg/id"
	"ymir.com/pkg/util"
	"ymir.com/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type SigninWithGoogleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSigninWithGoogleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SigninWithGoogleLogic {
	return &SigninWithGoogleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SigninWithGoogleLogic) SigninWithGoogle(req *types.SigninWithGoogleRequest) (*types.SigninResponse, error) {
	userGoogle, err := google.ValidateToken(l.ctx, req.Token)
	if err != nil {
		return nil, err
	}

	respb, err := l.svcCtx.UserRPC.GetUserGoogle(l.ctx, &user.GetUserGoogleRequest{
		GoogleUserId: userGoogle.GoogleUserId,
	})
	if err != nil {
		return nil, err
	}
	if respb.User != nil {
		return l.Signin(respb.User)
	}

	return l.SignupAndSignin(userGoogle)
}

func (l *SigninWithGoogleLogic) Signin(ugoogle *user.UserGoogleInfo) (resp *types.SigninResponse, err error) {
	respb, err := l.svcCtx.UserRPC.GetUserInfo(l.ctx, &user.GetUserInfoRequest{
		UserId:      &ugoogle.UserId,
	})
	if err != nil {
		return nil, err
	}
	if respb.User == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ServerCommonError), "cannot find user info by userId: %d", ugoogle.UserId)
	}

	var nowDate = time.Now().Unix()
	uIdEncoded, err := id.Hash.EncodedId(ugoogle.Id)
	if err != nil {
		return nil, err
	}
	token, err := util.GetJwtToken(l.svcCtx.Config.Auth.AccessSecret, nowDate, l.svcCtx.Config.Auth.AccessExpire, uIdEncoded)
	if err != nil {
		return nil, err
	}

	return &types.SigninResponse{
		UserId:      uIdEncoded,
		Username:    respb.User.Username,
		AccessToken: token,
		AvatarUrl:   respb.User.AvatarUrl,
	}, nil
}

func (l *SigninWithGoogleLogic) SignupAndSignin(usr *google.User) (resp *types.SigninResponse, err error) {
	respbInsert, err := l.svcCtx.UserRPC.WriteUserGoogleInDB(l.ctx, &user.WriteUserGoogleRequest{
		UserGoogleId: usr.GoogleUserId,
		Email:        usr.Email,
		AvatarUrl:    usr.AvatarUtl,
		UserName:     usr.UserName,
	})
	if err != nil {
		return nil, err
	}

	respbUser, err := l.svcCtx.UserRPC.GetUserInfo(l.ctx, &user.GetUserInfoRequest{
		UserId:      &respbInsert.UserId,
	})
	if err != nil {
		return nil, err
	}

	var nowDate = time.Now().Unix()
	uIdEncoded, err := id.Hash.EncodedId(respbUser.User.Id)
	if err != nil {
		return nil, err
	}
	token, err := util.GetJwtToken(l.svcCtx.Config.Auth.AccessSecret, nowDate, l.svcCtx.Config.Auth.AccessExpire, uIdEncoded)
	if err != nil {
		return nil, err
	}

	return &types.SigninResponse{
		UserId:      uIdEncoded,
		Username:    respbUser.User.Username,
		AccessToken: token,
		AvatarUrl:   respbUser.User.AvatarUrl,
	}, nil
}
