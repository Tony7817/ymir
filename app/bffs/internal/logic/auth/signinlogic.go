package auth

import (
	"context"
	"time"

	"ymir.com/app/bffs/internal/svc"
	"ymir.com/app/bffs/internal/types"
	useradmin "ymir.com/app/user/admin/user"
	"ymir.com/pkg/id"
	"ymir.com/pkg/util"
	"ymir.com/pkg/xerr"

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

func (l *SigninLogic) Signin(req *types.SigninRequest) (*types.SigninResponse, error) {
	respbOrg, err := l.svcCtx.UserAdminRPC.GetOrganizer(l.ctx, &useradmin.GetOrganizerRequest{
		Phonenumber: &req.Phonenumber,
	})
	if err != nil {
		return nil, err
	}

	if respbOrg.Organizer == nil {
		return nil, xerr.NewErrCode(xerr.UserNotExistedError)
	}

	// TODO: check captcha

	if req.Captcha != "928412" {
		return nil, xerr.NewErrCode(xerr.WrongCaptchaError)
	}

	var now = time.Now().Unix()
	token, err := util.GetJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, l.svcCtx.Config.Auth.AccessExpire, respbOrg.Organizer.Id)
	if err != nil {
		return nil, err
	}

	return &types.SigninResponse{
		UserId:      id.EncodeId(respbOrg.Organizer.Id),
		Name:        respbOrg.Organizer.Name,
		Phonenumber: respbOrg.Organizer.Phonenumber,
		AccessToken: token,
	}, nil
}
