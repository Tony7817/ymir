package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	ProductAdminRPC zrpc.RpcClientConf
	ProductRPC      zrpc.RpcClientConf
	StarAdminRPC    zrpc.RpcClientConf
	UserAdminRPC    zrpc.RpcClientConf
	Auth            JwtAuth
}

type JwtAuth struct {
	AccessSecret string
	AccessExpire int64
}
