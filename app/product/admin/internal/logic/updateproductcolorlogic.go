package logic

import (
	"context"

	"ymir.com/app/product/admin/internal/svc"
	"ymir.com/app/product/admin/product"
	"ymir.com/app/product/model"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductColorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProductColorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductColorLogic {
	return &UpdateProductColorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateProductColorLogic) UpdateProductColor(in *product.UpdateProductColorRequest) (*product.UpdateProductColorResponse, error) {
	ret, err := l.svcCtx.ProductColorModel.UpdatePartial(l.ctx, &model.ProductColorDetailPartial{
		Id:        in.Id,
		ProductId: in.ProductId,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "[UpdateProductColor] update product color fail")
	}

	uptRow, err := ret.RowsAffected()
	if err != nil {
		return nil, errors.Wrapf(err, "[UpdateProductColor] get affected row failed")
	}
	if uptRow == 0 {
		return nil, errors.Wrapf(errors.New("0 update"), "[UpdateProductColor] get affected row failed")
	}

	return &product.UpdateProductColorResponse{
		Id: in.Id,
	}, nil
}
