package logic

import (
	"context"
	"database/sql"

	"ymir.com/app/user/model"
	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/user"
	"ymir.com/pkg/id"
	"ymir.com/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type WriteUserLocalInDBLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWriteUserLocalInDBLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WriteUserLocalInDBLogic {
	return &WriteUserLocalInDBLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WriteUserLocalInDBLogic) WriteUserLocalInDB(in *user.WriteUserLocalRequest) (*user.WriteUserLocalResponse, error) {
	var usr = &model.User{
		Id:       id.SF.GenerateID(),
		Username: in.UserName,
	}
	if in.Email != nil {
		usr.Email = sql.NullString{
			String: *in.Email,
			Valid:  true,
		}
	} else if in.Phonenumber != nil {
		usr.PhoneNumber = sql.NullString{
			String: *in.Phonenumber,
			Valid:  true,
		}
	} else {
		return nil, xerr.NewErrCode(xerr.ErrorReuqestParam)
	}

	var userId = id.SF.GenerateID()
	err := l.svcCtx.UserModel.InsertIntoUserAndUserLocal(l.ctx, usr, &model.UserLocal{
		Id:           userId,
		UserId:       usr.Id,
		PasswordHash: in.PasswordHash,
		IsActivated:  1,
	})
	if err != nil {
		return nil, err
	}

	return &user.WriteUserLocalResponse{
		UserId: userId,
	}, nil
}
