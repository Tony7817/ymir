package product

import (
	"context"

	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/app/product/rpc/product"
	"ymir.com/pkg/id"

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
	var (
		pIdDecoded = id.Hash.DecodedId(req.ProductId)
		cIdDecoded = id.Hash.DecodedId(req.ColorId)
	)

	respb, err := l.svcCtx.ProductRPC.AddProductToCart(l.ctx, &product.AddProductToCartRequest{
		ProductId: pIdDecoded,
		UserId:    uIdDecoded,
		Size:      req.Size,
		ColorId:   cIdDecoded,
	})
	if err != nil {
		return nil, err
	}

	pcIdEncoded, err := id.Hash.EncodedId(respb.ProductCartId)
	if err != nil {
		return nil, err
	}

	return &types.AddProductToCartResponse{
		ProductCartId: pcIdEncoded,
	}, nil
}
