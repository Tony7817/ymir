package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ProductCartModel = (*customProductCartModel)(nil)

type (
	// ProductCartModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductCartModel.
	ProductCartModel interface {
		productCartModel
		FindProductsOfUser(userId int64, page int64, pageSize int64) ([]*ProductCart, error)
		CountTotalProductOfUser(ctx context.Context, userId int64) (int64, error)
	}

	customProductCartModel struct {
		*defaultProductCartModel
	}
)

// NewProductCartModel returns a model for the database table.
func NewProductCartModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ProductCartModel {
	return &customProductCartModel{
		defaultProductCartModel: newProductCartModel(conn, c, opts...),
	}
}

func (m *customProductCartModel) FindProductsOfUser(userId int64, page int64, pageSize int64) ([]*ProductCart, error) {
	if pageSize > 10 {
		pageSize = 10
	}

	var res []*ProductCart
	var offset = (page - 1) * pageSize
	if err := m.QueryRowsNoCache(&res, "SELECT * FROM ? WHERE user_id = ? LIMIT ?, ?", m.table, userId, offset, pageSize); err != nil {
		return nil, err
	}

	return res, nil
}

func (m *customProductCartModel) CountTotalProductOfUser(ctx context.Context, userId int64) (int64, error) {
	var count int64
	var cacheKey = "cache:productCart:userId:total"
	err := m.QueryRowCtx(ctx, &count, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		return conn.QueryRow(ctx, "SELECT COUNT(*) FROM ? WHERE user_id = ?", m.table, userId)
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}
