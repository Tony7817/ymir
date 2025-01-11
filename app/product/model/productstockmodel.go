package model

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ymir.com/pkg/id"
)

var _ ProductStockModel = (*customProductStockModel)(nil)

type (
	// ProductStockModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductStockModel.
	ProductStockModel interface {
		productStockModel
		InsertSF(ctx context.Context, data *ProductStock) (int64, error)
	}

	customProductStockModel struct {
		*defaultProductStockModel
		sf *id.Snowflake
	}
)

// NewProductStockModel returns a model for the database table.
func NewProductStockModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ProductStockModel {
	return &customProductStockModel{
		defaultProductStockModel: newProductStockModel(conn, c, opts...),
		sf:                       id.NewSnowFlake(),
	}
}

func (m *customProductStockModel) InsertSF(ctx context.Context, data *ProductStock) (int64, error) {
	data.Id = m.sf.GenerateID()
	if _, err := m.Insert(ctx, data); err != nil {
		return 0, errors.Wrapf(err, "[InsertSF]insert product stock fail")
	}

	return data.Id, nil
}
