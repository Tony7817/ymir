package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PaypalModel = (*customPaypalModel)(nil)

type (
	// PaypalModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPaypalModel.
	PaypalModel interface {
		paypalModel
	}

	customPaypalModel struct {
		*defaultPaypalModel
	}
)

// NewPaypalModel returns a model for the database table.
func NewPaypalModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PaypalModel {
	return &customPaypalModel{
		defaultPaypalModel: newPaypalModel(conn, c, opts...),
	}
}
