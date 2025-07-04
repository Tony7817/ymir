package main

import (
	"flag"
	"fmt"
	"os"

	"ymir.com/app/order/rpc/internal/config"
	"ymir.com/app/order/rpc/internal/server"
	"ymir.com/app/order/rpc/internal/svc"
	"ymir.com/app/order/rpc/order"
	"ymir.com/pkg/interceptor"
	"ymir.com/pkg/paypal"
	"ymir.com/pkg/vars"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/order.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	if c.Mode == vars.ModeProd {
		os.Setenv(vars.ModeVar, vars.ModeProd)
	} else if c.Mode == vars.ModeDev {
		os.Setenv(vars.ModeVar, vars.ModeDev)
	}

	paypal.Init()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		order.RegisterOrderServer(grpcServer, server.NewOrderServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	s.AddUnaryInterceptors(interceptor.LoggerInterceptor)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
