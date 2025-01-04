package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"ymir.com/app/bffs/internal/config"
	"ymir.com/app/bffs/internal/middleware"
	"ymir.com/app/star/admin/starclient"
	useradminclient "ymir.com/app/user/admin/userclient"
	"ymir.com/app/user/rpc/userclient"
)

type ServiceContext struct {
	Config       config.Config
	StarRPC      starclient.Star
	UserAdminRPC useradminclient.User
	UserRPC      userclient.User
	Auth         rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		StarRPC:      starclient.NewStar(zrpc.MustNewClient(c.StarRPC)),
		UserAdminRPC: useradminclient.NewUser(zrpc.MustNewClient(c.UserRPC)),
		UserRPC:      userclient.NewUser(zrpc.MustNewClient(c.UserRPC)),
		Auth:         middleware.NewAuthMiddleware(c).Handle,
	}
}
