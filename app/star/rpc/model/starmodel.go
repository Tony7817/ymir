package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ StarModel = (*customStarModel)(nil)

type (
	// StarModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStarModel.
	StarModel interface {
		starModel
		FindStarList(ctx context.Context, offset, limit int64) ([]*Star, error)
		CountStarTotal(ctx context.Context) (int64, error)
	}

	customStarModel struct {
		*defaultStarModel
	}
)

// NewStarModel returns a model for the database table.
func NewStarModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) StarModel {
	return &customStarModel{
		defaultStarModel: newStarModel(conn, c, opts...),
	}
}

func (m *customStarModel) FindStarList(ctx context.Context, offset, limit int64) ([]*Star, error) {
	var stars []*Star
	err := m.QueryRowsNoCacheCtx(ctx, &stars, "SELECT * FROM `star` LIMIT ?,?", offset, limit)
	if err != nil {
		return nil, err
	}

	return stars, nil
}

func (m *customStarModel) CountStarTotal(ctx context.Context) (int64, error) {
	const cacheKey = "cache:star:total"
	var total int64
	if err := m.QueryRowCtx(ctx, &total, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		return conn.QueryRowCtx(ctx, &total, "SELECT COUNT(*) FROM `star`")
	}); err != nil {
		return 0, err
	}

	return total, nil
}
