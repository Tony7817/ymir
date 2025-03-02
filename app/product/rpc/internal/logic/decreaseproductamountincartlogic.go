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

type DecreaseProductAmountInCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDecreaseProductAmountInCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DecreaseProductAmountInCartLogic {
	return &DecreaseProductAmountInCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DecreaseProductAmountInCartLogic) DecreaseProductAmountInCart(in *product.DecreaseProductAmountInCartRequest) (*product.AddProductAmountInCartResponse, error) {
	pc, err := l.svcCtx.ProductCartModel.FindOneByUserIdColorIdProductIdSize(l.ctx, in.UserId, in.ColorId, in.ProductId, in.Size)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "find product cart failed")
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.DataNoExistError), "[DecreaseProductAmountInCart]product cart not exist")
	}

	if pc.Amount == 0 {
		return &product.AddProductAmountInCartResponse{}, nil
	}

	if err := l.svcCtx.ProductCartModel.DescProductAmount(l.ctx, in.UserId, in.ProductId, in.ColorId, in.Size); err != nil {
		return nil, err
	}

	var pcNew *model.ProductCart
	var pcd *model.ProductColorDetail
	err = mr.Finish(func() error {
		var err error
		pcNew, err = l.svcCtx.ProductCartModel.FindOneByUserIdColorIdProductIdSize(l.ctx, in.UserId, in.ColorId, in.ProductId, in.Size)
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
		ProductCartId: pcNew.Id,
		Amount:        pcNew.Amount,
		TotalPrice:    pcNew.Amount * pcd.Price,
	}, nil
}
