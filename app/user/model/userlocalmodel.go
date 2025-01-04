package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserLocalModel = (*customUserLocalModel)(nil)

type (
	// UserLocalModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserLocalModel.
	UserLocalModel interface {
		userLocalModel
	}

	customUserLocalModel struct {
		*defaultUserLocalModel
	}
)

// NewUserLocalModel returns a model for the database table.
func NewUserLocalModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserLocalModel {
	return &customUserLocalModel{
		defaultUserLocalModel: newUserLocalModel(conn, c, opts...),
	}
}
