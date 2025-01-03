package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	OrganizerRPC zrpc.RpcClientConf
	ProductRPC   zrpc.RpcClientConf
	StarRPC      zrpc.RpcClientConf
	UserRPC      zrpc.RpcClientConf
	Auth         JwtAuth
}

type JwtAuth struct {
	AccessSecret string
	AccessExpire int64
}
