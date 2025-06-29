// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.3
// Source: useradmin.proto

package server

import (
	"context"

	"ymir.com/app/user/admin/internal/logic"
	"ymir.com/app/user/admin/internal/svc"
	"ymir.com/app/user/admin/user"
)

type UserServer struct {
	svcCtx *svc.ServiceContext
	user.UnimplementedUserServer
}

func NewUserServer(svcCtx *svc.ServiceContext) *UserServer {
	return &UserServer{
		svcCtx: svcCtx,
	}
}

func (s *UserServer) GetOrganizer(ctx context.Context, in *user.GetOrganizerRequest) (*user.GetOrganizerResponse, error) {
	l := logic.NewGetOrganizerLogic(ctx, s.svcCtx)
	return l.GetOrganizer(in)
}
