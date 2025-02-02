package star

import (
	"context"

	"ymir.com/app/bffs/internal/svc"
	"ymir.com/app/bffs/internal/types"
	"ymir.com/app/product/admin/product"
	"ymir.com/pkg/id"
	"ymir.com/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProductLogic {
	return &DeleteProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteProductLogic) DeleteProduct(req *types.DeleteProductRequest) (resp *types.DeleteProductResponse, err error) {
	pId, err := id.DecodeId(req.ProductId)
	if err != nil {
		return nil, xerr.NewErrCode(xerr.ErrorNotAuthorized)
	}

	respb, err := l.svcCtx.ProductAdminRPC.DeleteProduct(l.ctx, &product.DeleteProductRequest{
		ProductId: pId,
	})
	if err != nil {
		return nil, err
	}

	return &types.DeleteProductResponse{
		ProductId: id.EncodeId(respb.ProductId),
	}, nil
}
