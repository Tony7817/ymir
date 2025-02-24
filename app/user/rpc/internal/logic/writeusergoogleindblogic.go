package logic

import (
	"context"
	"database/sql"

	"ymir.com/app/user/model"
	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/user"
	"ymir.com/pkg/id"

	"github.com/zeromicro/go-zero/core/logx"
)

type WriteUserGoogleInDBLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWriteUserGoogleInDBLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WriteUserGoogleInDBLogic {
	return &WriteUserGoogleInDBLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WriteUserGoogleInDBLogic) WriteUserGoogleInDB(in *user.WriteUserGoogleRequest) (*user.WriteUserGoogleResponse, error) {
	var userId = id.SF.GenerateID()
	err := l.svcCtx.UserModel.InsertIntoUserAndUserGoogle(l.ctx, &model.User{
		Id:       userId,
		Username: in.UserName,
		Email: sql.NullString{
			String: in.Email,
			Valid:  true,
		},
		AvatarUrl: sql.NullString{
			String: in.AvatarUrl,
			Valid:  true,
		},
	}, &model.UserGoogle{
		GoogleUserId: in.UserGoogleId,
	})
	if err != nil {
		return nil, err
	}

	return &user.WriteUserGoogleResponse{
		UserId: userId,
	}, nil
}
