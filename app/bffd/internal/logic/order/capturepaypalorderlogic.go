package order

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CapturePaypalOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCapturePaypalOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CapturePaypalOrderLogic {
	return &CapturePaypalOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CapturePaypalOrderLogic) CapturePaypalOrder(req *types.CapturePaypalOrderRequest) (resp *types.CapturePaypalOrderResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
