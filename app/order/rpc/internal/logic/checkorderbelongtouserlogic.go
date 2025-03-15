package logic

import (
	"context"
	"database/sql"

	"ymir.com/app/order/rpc/internal/svc"
	"ymir.com/app/order/rpc/order"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CheckOrderBelongToUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckOrderBelongToUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckOrderBelongToUserLogic {
	return &CheckOrderBelongToUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckOrderBelongToUserLogic) CheckOrderBelongToUser(in *order.CheckOrderBelongTouserRequest) (*order.CheckOrderBelongTouserResponse, error) {
	_, err := l.svcCtx.OrderModel.FindOneByIdUserId(l.ctx, in.OrderId, in.UserId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrapf(err, "find order by id and user id failed, orderId: %d, userId: %d", in.OrderId, in.UserId)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return &order.CheckOrderBelongTouserResponse{
			IsBelong: false,
		}, nil
	}

	return &order.CheckOrderBelongTouserResponse{
		IsBelong: true,
	}, nil
}
