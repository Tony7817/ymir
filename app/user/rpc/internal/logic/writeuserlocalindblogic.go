package logic

import (
	"context"
	"database/sql"

	"ymir.com/app/user/model"
	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/user"
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
		Username: in.UserName,
		AvatarUrl: sql.NullString{
			String: *in.AvatarUrl,
			Valid:  true,
		},
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
		return nil, xerr.NewErrCode(xerr.ReuqestParamError)
	}

	uId, err := l.svcCtx.UserModel.InsertIntoUserAndUserLocal(l.ctx, usr, &model.UserLocal{
		PasswordHash: in.PasswordHash,
		IsActivated:  1,
	})
	if err != nil {
		return nil, err
	}

	return &user.WriteUserLocalResponse{
		UserId: uId,
	}, nil
}
