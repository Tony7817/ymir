package svc

import (
	"ymir.com/app/bffd/internal/config"
	"ymir.com/app/bffd/internal/middleware"
	"ymir.com/app/product/rpc/productclient"
	"ymir.com/app/star/rpc/starclient"
	"ymir.com/app/user/rpc/userclient"
	"ymir.com/pkg/id"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	Hash       *id.HashID
	ProductRPC productclient.Product
	StarRPC    starclient.Star
	UserRPC    userclient.User

	Timer rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		Hash:       id.NewHashID(),
		ProductRPC: productclient.NewProduct(zrpc.MustNewClient(c.ProductRPC)),
		StarRPC:    starclient.NewStar(zrpc.MustNewClient(c.StarRPC)),
		UserRPC:    userclient.NewUser(zrpc.MustNewClient(c.UserRPC)),
		//
		Timer: middleware.NewTimerMiddleware(c).Handle,
	}
}
