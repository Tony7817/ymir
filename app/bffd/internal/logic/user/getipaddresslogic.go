package user

import (
	"context"
	"net"
	"net/http"
	"strings"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetIpAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetIpAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIpAddressLogic {
	return &GetIpAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetIpAddressLogic) GetIpAddress(req *http.Request) (resp *types.GetIpAddressResponse, err error) {
	var res = &types.GetIpAddressResponse{}
	if ip := req.Header.Get("X-Forwarded-For"); ip != "" {
		ips := strings.Split(ip, ",")
		res.Ip = ips[0]
		return res, nil
	}

	if ip := req.Header.Get("X-Real-Ip"); ip != "" {
		res.Ip = ip
		return res, nil
	}

	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		res.Ip = req.RemoteAddr
	} else {
		res.Ip = ip
		return res, nil
	}

	return nil, xerr.NewErrCode(xerr.ReuqestParamError)
}
