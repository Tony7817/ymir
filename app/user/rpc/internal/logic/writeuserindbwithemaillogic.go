package logic

import (
	"context"
	"database/sql"

	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/model"
	"ymir.com/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type WriteUserInDBWithEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWriteUserInDBWithEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WriteUserInDBWithEmailLogic {
	return &WriteUserInDBWithEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WriteUserInDBWithEmailLogic) WriteUserInDBWithEmail(in *user.WriteUserInDBWithEmailRequest) (*user.WriteUserInDBWithEmailResponse, error) {
	var usr = &model.User{
		Username: in.UserName,
		Email: sql.NullString{
			String: in.Email,
			Valid:  true,
		},
		PasswordHash: in.PasswordHash,
		IsActivated:  1,
	}
	ret, err := l.svcCtx.UserModel.Insert(l.ctx, usr)
	if err != nil {
		return nil, err
	}

	lastId, err := ret.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &user.WriteUserInDBWithEmailResponse{
		UserId: lastId,
	}, nil
}
