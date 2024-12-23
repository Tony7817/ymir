package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ProductModel = (*customProductModel)(nil)

type (
	// ProductModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductModel.
	ProductModel interface {
		productModel
		FindProductList(ctx context.Context, starId *int64, offset, limit int64) ([]*Product, error)
		CountTotalProduct(ctx context.Context, starId *int64) (int64, error)
		CheckProductBelongToStar(ctx context.Context, productId, starId int64) (bool, error)
	}

	customProductModel struct {
		*defaultProductModel
	}
)

// NewProductModel returns a model for the database table.
func NewProductModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ProductModel {
	return &customProductModel{
		defaultProductModel: newProductModel(conn, c, opts...),
	}
}

func (m *customProductModel) CheckProductBelongToStar(ctx context.Context, productId, starId int64) (bool, error) {
	var count int64
	var cacheKey = fmt.Sprintf("cache:starId:of:product:%d", productId)
	err := m.QueryRowCtx(ctx, &count, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		return conn.QueryRowCtx(ctx, &count, "select count(*) from product where id = ? and star_id = ?", productId, starId)
	})
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func (m *customProductModel) FindProductList(ctx context.Context, starId *int64, offset, limit int64) ([]*Product, error) {
	if limit > 20 {
		limit = 20
	}
	var ps []*Product

	var cond = " where 1=1"
	if starId != nil {
		cond += fmt.Sprintf(" and star_id = %d", *starId)
	}
	cond += " limit ?,?"

	var query = fmt.Sprintf("SELECT * FROM `product` %s", cond)
	err := m.QueryRowsNoCacheCtx(ctx, &ps, query, offset, limit)
	if err != nil {
		return nil, err
	}

	return ps, nil
}

func (m *customProductModel) CountTotalProduct(ctx context.Context, starId *int64) (int64, error) {
	var cacheKey string
	if starId != nil {
		cacheKey = fmt.Sprintf("cache:total:product:star:%d", *starId)
	} else {
		cacheKey = "cache:total:product"
	}
	var count int64
	var cond = " where 1=1"
	if starId != nil {
		cond += fmt.Sprintf(" and star_id = %d", *starId)
	}
	err := m.QueryRowCtx(ctx, &count, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		return conn.QueryRowCtx(ctx, &count, "select count(*) from product %s", cond)
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}
