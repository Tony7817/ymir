package logic

import (
	"context"
	"database/sql"

	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/model"
	"ymir.com/app/user/rpc/user"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserGoogleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserGoogleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserGoogleLogic {
	return &GetUserGoogleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserGoogleLogic) GetUserGoogle(in *user.GetUserGoogleRequest) (*user.GetUserGoogleResponse, error) {
	var (
		usrGoogle *model.UserGoogle
		err       error
	)

	usrGoogle, err = l.svcCtx.UserGoogleModel.FindOneByGoogleUserId(l.ctx, in.GoogleUserId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return &user.GetUserGoogleResponse{User: nil}, nil
	}

	return &user.GetUserGoogleResponse{
		User: &user.UserGoogleInfo{
			Id:           usrGoogle.Id,
			UserId:       usrGoogle.UserId,
			GoogleUserId: usrGoogle.GoogleUserId,
		},
	}, nil
}
