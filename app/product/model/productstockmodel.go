package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ProductStockModel = (*customProductStockModel)(nil)

type (
	// ProductStockModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductStockModel.
	ProductStockModel interface {
		productStockModel
	}

	customProductStockModel struct {
		*defaultProductStockModel
	}
)

// NewProductStockModel returns a model for the database table.
func NewProductStockModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ProductStockModel {
	return &customProductStockModel{
		defaultProductStockModel: newProductStockModel(conn, c, opts...),
	}
}
