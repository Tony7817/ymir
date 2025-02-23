package product

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/product/rpc/product"
	"ymir.com/pkg/id"
	"ymir.com/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddProductToCartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddProductToCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddProductToCartLogic {
	return &AddProductToCartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddProductToCartLogic) AddProductToCart(req *types.AddProductToCartRequest) (resp *types.AddProductToCartResponse, err error) {
	uIdDecoded, err := id.GetDecodedUserId(l.ctx)
	if err != nil {
		return nil, err
	}
	pId, err := id.DecodeId(req.ProductId)
	if err != nil {
		return nil, err
	}
	cId, err := id.DecodeId(req.ColorId)
	if err != nil {
		return nil, err
	}

	stock, err := l.svcCtx.ProductRPC.ProductStock(l.ctx, &product.ProductStockRequest{
		ProductId: pId,
		ColorId:   cId,
		Size:      req.Size,
	})
	if err != nil {
		return nil, err
	}
	if stock.Stock <= 0 {
		return nil, xerr.NewErrCode(xerr.ErrorOutOfStockError)
	}

	respb, err := l.svcCtx.ProductRPC.AddProductToCart(l.ctx, &product.AddProductToCartRequest{
		ProductId: pId,
		UserId:    uIdDecoded,
		Size:      req.Size,
		ColorId:   cId,
	})
	if err != nil {
		return nil, err
	}

	return &types.AddProductToCartResponse{
		ProductCartId: id.EncodeId(respb.ProductCartId),
	}, nil
}
