package logic

import (
	"context"
	"database/sql"

	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/model"
	"ymir.com/app/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type WriteUserInDBWithPhonenumberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWriteUserInDBWithPhonenumberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WriteUserInDBWithPhonenumberLogic {
	return &WriteUserInDBWithPhonenumberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WriteUserInDBWithPhonenumberLogic) WriteUserInDBWithPhonenumber(in *user.WriteUserInDBWithPhonenumberRequest) (*user.WriteUserInDBWithPhonenumberResponse, error) {
	var usr = &model.User{
		Username: in.UserName,
		PhoneNumber: sql.NullString{
			String: in.Phonenumber,
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

	return &user.WriteUserInDBWithPhonenumberResponse{
		UserId: lastId,
	}, nil
}
