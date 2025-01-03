package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"ymir.com/app/bffs/internal/config"
	"ymir.com/app/organizer/rpc/organizerclient"
	"ymir.com/app/product/rpc/productclient"
	"ymir.com/app/star/rpc/starclient"
	"ymir.com/app/user/rpc/userclient"
)

type ServiceContext struct {
	Config       config.Config
	OrganizerRPC organizerclient.Organizer
	ProductRPC   productclient.Product
	StarRPC      starclient.Star
	UserRPC      userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		OrganizerRPC: organizerclient.NewOrganizer(zrpc.MustNewClient(c.OrganizerRPC)),
		ProductRPC:   productclient.NewProduct(zrpc.MustNewClient(c.ProductRPC)),
		StarRPC:      starclient.NewStar(zrpc.MustNewClient(c.StarRPC)),
		UserRPC:      userclient.NewUser(zrpc.MustNewClient(c.UserRPC)),
	}
}
