package logic

import (
	"context"
	"database/sql"

	"ymir.com/app/order/rpc/internal/svc"
	"ymir.com/app/order/rpc/order"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CheckOrderIdempotentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckOrderIdempotentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckOrderIdempotentLogic {
	return &CheckOrderIdempotentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckOrderIdempotentLogic) CheckOrderIdempotent(in *order.CheckOrderIdempotentRequest) (*order.CheckOrderIdempotentResponse, error) {
	_, err := l.svcCtx.OrderModel.FindOneByRequestId(l.ctx, in.RequestId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "find order by request id failed")
	}
	if errors.Is(err, sql.ErrNoRows) {
		return &order.CheckOrderIdempotentResponse{
			IsIdempotent: true,
		}, nil
	}

	return &order.CheckOrderIdempotentResponse{
		IsIdempotent: false,
	}, nil
}
