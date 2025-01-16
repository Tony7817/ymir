package logic

import (
	"context"
	"database/sql"

	"ymir.com/app/product/admin/internal/svc"
	"ymir.com/app/product/admin/product"
	"ymir.com/app/product/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductLogic {
	return &CreateProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateProductLogic) CreateProduct(in *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	var desc sql.NullString
	if in.Detail != nil && *in.Detail != "" {
		desc.String = *in.Detail
		desc.Valid = true
	}

	id, err := l.svcCtx.ProductModel.SFInsert(l.ctx, &model.Product{
		StarId:         in.StarId,
		Name:           in.Name,
		SoldNum:        0,
		Description:    in.Description,
		Rate:           5,
		RateCount:      0,
		Detail: desc,
	})
	if err != nil {
		return nil, err
	}

	return &product.CreateProductResponse{
		Id: id,
	}, nil
}
