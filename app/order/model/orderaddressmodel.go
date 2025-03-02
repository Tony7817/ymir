package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderAddressModel = (*customOrderAddressModel)(nil)

type (
	// OrderAddressModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderAddressModel.
	OrderAddressModel interface {
		orderAddressModel
	}

	customOrderAddressModel struct {
		*defaultOrderAddressModel
	}
)

// NewOrderAddressModel returns a model for the database table.
func NewOrderAddressModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrderAddressModel {
	return &customOrderAddressModel{
		defaultOrderAddressModel: newOrderAddressModel(conn, c, opts...),
	}
}
