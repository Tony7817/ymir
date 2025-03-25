package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"ymir.com/app/bffd/internal/config"
	"ymir.com/app/bffd/internal/handler"
	"ymir.com/app/bffd/internal/svc"
	"ymir.com/pkg/vars"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/bff.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	if c.Mode == vars.ModeProd {
		os.Setenv(vars.ModeVar, vars.ModeProd)
	} else if c.Mode == vars.ModeDev {
		os.Setenv(vars.ModeVar, vars.ModeDev)
	}

	server := rest.MustNewServer(c.RestConf, rest.WithCors("*"), rest.WithCorsHeaders("X-Content-Security"), rest.WithUnsignedCallback(func(w http.ResponseWriter, r *http.Request, next http.Handler, strict bool, code int) {
		http.Error(w, fmt.Sprintf("unsafe request, code:%d, strict:%+v", code, strict), http.StatusForbidden)
	}))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
