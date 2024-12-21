package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	ProductRPC zrpc.RpcClientConf
	StarRPC    zrpc.RpcClientConf
	UserRPC    zrpc.RpcClientConf
	Auth       JwtAuth

	BizRedis redis.RedisConf
}

type JwtAuth struct {
	AccessSecret string
	AccessExpire int64
}
