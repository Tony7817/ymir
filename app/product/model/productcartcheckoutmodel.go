package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ProductCartCheckoutModel = (*customProductCartCheckoutModel)(nil)

type (
	// ProductCartCheckoutModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductCartCheckoutModel.
	ProductCartCheckoutModel interface {
		productCartCheckoutModel
	}

	customProductCartCheckoutModel struct {
		*defaultProductCartCheckoutModel
	}
)

// NewProductCartCheckoutModel returns a model for the database table.
func NewProductCartCheckoutModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ProductCartCheckoutModel {
	return &customProductCartCheckoutModel{
		defaultProductCartCheckoutModel: newProductCartCheckoutModel(conn, c, opts...),
	}
}
