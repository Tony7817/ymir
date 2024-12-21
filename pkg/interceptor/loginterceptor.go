package interceptor

import (
	"context"

	"ymir.com/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LoggerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		causeErr := errors.Cause(err)
		if e, ok := causeErr.(*xerr.CodeError); ok {
			logx.WithContext(ctx).Errorf("rpc-srv-err %+v", err)
			err = status.Error(codes.Code(e.GetErrCode()), e.GetErrMsg())
		} else {
			logx.WithContext(ctx).Errorf("rpc-srv-err %+v", err)
		}
	}

	return resp, err
}
