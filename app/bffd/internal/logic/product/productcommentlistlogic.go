package product

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductCommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProductCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductCommentListLogic {
	return &ProductCommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProductCommentListLogic) ProductCommentList(req *types.ProductCommentRequest) (resp *types.ProductCommentResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
