package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"ymir.com/app/bffs/internal/config"
	"ymir.com/app/bffs/internal/middleware"
	productadminclient "ymir.com/app/product/admin/productclient"
	"ymir.com/app/product/rpc/productclient"
	"ymir.com/app/star/admin/starclient"
	useradminclient "ymir.com/app/user/admin/userclient"
)

type ServiceContext struct {
	Config          config.Config
	StarAdminRPC    starclient.Star
	UserAdminRPC    useradminclient.User
	ProductAdminRPC productadminclient.Product
	ProductRPC      productclient.Product
	Auth            rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:          c,
		StarAdminRPC:    starclient.NewStar(zrpc.MustNewClient(c.StarAdminRPC)),
		UserAdminRPC:    useradminclient.NewUser(zrpc.MustNewClient(c.UserAdminRPC)),
		ProductRPC:      productclient.NewProduct(zrpc.MustNewClient(c.ProductRPC)),
		ProductAdminRPC: productadminclient.NewProduct(zrpc.MustNewClient(c.ProductAdminRPC)),
		Auth:            middleware.NewAuthMiddleware(c).Handle,
	}
}
