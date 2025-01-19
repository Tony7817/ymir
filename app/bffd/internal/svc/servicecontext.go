package svc

import (
	"ymir.com/app/bffd/internal/config"
	"ymir.com/app/bffd/internal/middleware"
	"ymir.com/app/order/rpc/orderclient"
	"ymir.com/app/product/rpc/productclient"
	"ymir.com/app/star/rpc/starclient"
	"ymir.com/app/user/rpc/userclient"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	ProductRPC productclient.Product
	StarRPC    starclient.Star
	UserRPC    userclient.User
	OrderRPC   orderclient.Order

	Timer rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		ProductRPC: productclient.NewProduct(zrpc.MustNewClient(c.ProductRPC)),
		StarRPC:    starclient.NewStar(zrpc.MustNewClient(c.StarRPC)),
		UserRPC:    userclient.NewUser(zrpc.MustNewClient(c.UserRPC)),
		OrderRPC:   orderclient.NewOrder(zrpc.MustNewClient(c.OrderRPC)),
		//
		Timer: middleware.NewTimerMiddleware(c).Handle,
	}
}
