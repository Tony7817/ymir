package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrganizerModel = (*customOrganizerModel)(nil)

type (
	// OrganizerModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrganizerModel.
	OrganizerModel interface {
		organizerModel
	}

	customOrganizerModel struct {
		*defaultOrganizerModel
	}
)

// NewOrganizerModel returns a model for the database table.
func NewOrganizerModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrganizerModel {
	return &customOrganizerModel{
		defaultOrganizerModel: newOrganizerModel(conn, c, opts...),
	}
}
