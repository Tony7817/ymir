package logic

import (
	"context"
	"database/sql"

	"ymir.com/app/product/model"
	"ymir.com/app/product/rpc/internal/svc"
	"ymir.com/app/product/rpc/product"
	"ymir.com/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type AddProductAmountInCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddProductAmountInCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddProductAmountInCartLogic {
	return &AddProductAmountInCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddProductAmountInCartLogic) AddProductAmountInCart(in *product.AddProductAmountInCartRequest) (*product.AddProductAmountInCartResponse, error) {
	ps, err := l.svcCtx.ProductStockModel.FindOneByProductIdColorIdSize(l.ctx, in.ProductId, in.ColorId, in.Size)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, xerr.NewErrCode(xerr.DataNoExistError)
	}
	if ps.InStock < in.ExpectedAmount {
		return nil, xerr.NewErrCode(xerr.ErrorStockNotEnough)
	}

	if err := l.svcCtx.ProductCartModel.IncrProductAmount(l.ctx, in.UserId, in.ProductId, in.ColorId, in.Size); err != nil {
		return nil, err
	}

	var pc *model.ProductCart
	var pcd *model.ProductColorDetail
	err = mr.Finish(func() error {
		var err error
		pc, err = l.svcCtx.ProductCartModel.FindOneByUserIdColorIdProductIdSize(l.ctx, in.UserId, in.ColorId, in.ProductId, in.Size)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errors.Wrap(err, "find product cart failed")
		}

		if errors.Is(err, sql.ErrNoRows) {
			return xerr.NewErrCode(xerr.DataNoExistError)
		}
		return nil
	}, func() error {
		var err error
		pcd, err = l.svcCtx.ProductColorModel.FindOne(l.ctx, in.ColorId)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errors.Wrap(err, "find product failed")
		}
		if errors.Is(err, sql.ErrNoRows) {
			return xerr.NewErrCode(xerr.DataNoExistError)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &product.AddProductAmountInCartResponse{
		ProductCartId: pc.Id,
		Amount:        pc.Amount,
		TotalPrice:    pc.Amount * pcd.Price,
	}, nil
}
