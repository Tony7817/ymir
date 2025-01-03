package logic

import (
	"context"

	"ymir.com/app/organizer/rpc/internal/svc"
	"ymir.com/app/organizer/rpc/organizer"

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

func (l *GetOrganizerLogic) GetOrganizer(in *organizer.GetOrganizerRequest) (*organizer.GetorganizerResponse, error) {
	// todo: add your logic here and delete this line

	return &organizer.GetorganizerResponse{}, nil
}
