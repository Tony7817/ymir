package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserGoogleModel = (*customUserGoogleModel)(nil)

type (
	// UserGoogleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserGoogleModel.
	UserGoogleModel interface {
		userGoogleModel
	}

	customUserGoogleModel struct {
		*defaultUserGoogleModel
	}
)

// NewUserGoogleModel returns a model for the database table.
func NewUserGoogleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserGoogleModel {
	return &customUserGoogleModel{
		defaultUserGoogleModel: newUserGoogleModel(conn, c, opts...),
	}
}
