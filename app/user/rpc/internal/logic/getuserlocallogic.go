package logic

import (
	"context"
	"database/sql"

	"ymir.com/app/user/model"
	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/user"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLocalLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLocalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLocalLogic {
	return &GetUserLocalLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLocalLogic) GetUserLocal(in *user.GetUserLocalRequest) (*user.GetUserLocalResponse, error) {
	var (
		usr *model.UserLocal
		err error
	)

	usr, err = l.svcCtx.UserLocalModel.FindOneByUserId(l.ctx, in.UserId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return &user.GetUserLocalResponse{
			User: nil,
		}, nil
	}

	var res = &user.GetUserLocalResponse{
		User: &user.UserLocalInfo{
			Id:           usr.Id,
			UserId:       usr.UserId,
			PasswordHash: usr.PasswordHash,
			IsActivated:  usr.IsActivated == 1,
		},
	}

	return res, nil
}
