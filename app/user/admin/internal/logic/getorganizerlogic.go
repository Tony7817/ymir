package logic

import (
	"context"
	"database/sql"
	"errors"

	"ymir.com/app/user/admin/internal/svc"
	"ymir.com/app/user/admin/user"
	"ymir.com/app/user/model"
	"ymir.com/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrganizerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrganizerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrganizerLogic {
	return &GetOrganizerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOrganizerLogic) GetOrganizer(in *user.GetOrganizerRequest) (*user.GetOrganizerResponse, error) {
	var (
		o   *model.Organizer
		err error
	)
	if in.Phonenumber != nil {
		o, err = l.svcCtx.OrganizerModel.FindOneByPhoneNumber(l.ctx, sql.NullString{
			String: *in.Phonenumber,
			Valid:  true,
		})
	} else if in.UserId != nil {
		o, err = l.svcCtx.OrganizerModel.FindOne(l.ctx, *in.UserId)
	} else {
		return nil, xerr.NewErrCode(xerr.ReuqestParamError)
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return &user.GetOrganizerResponse{
			Organizer: nil,
		}, nil
	}

	return &user.GetOrganizerResponse{
		Organizer: &user.Organizer{
			Id:          o.Id,
			Phonenumber: o.PhoneNumber.String,
			Role:        o.Role,
		},
	}, nil
}
