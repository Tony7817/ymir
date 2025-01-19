package model

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/threading"
	"ymir.com/pkg/util"
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
		SFInsert(ctx context.Context, product *Product) (int64, error)
		DeleteProduct(ctx context.Context, pId int64) error
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
	var cacheKey = CacheKeyProductBelongToStar(productId)
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

	var query = fmt.Sprintf("select * from `product` %s", cond)
	err := m.QueryRowsNoCacheCtx(ctx, &ps, query, offset, limit)
	if err != nil {
		return nil, errors.Wrapf(err, "[FindProductList]query products list failed")
	}

	return ps, nil
}

func (m *customProductModel) CountTotalProduct(ctx context.Context, starId *int64) (int64, error) {
	var cacheKey = CacheKeyProductCount(starId)
	var count int64
	var cond = " where 1=1"
	if starId != nil {
		cond += fmt.Sprintf(" and star_id = %d", *starId)
	}
	err := m.QueryRowCtx(ctx, &count, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		return conn.QueryRowCtx(ctx, &count, "select count(*) from product"+cond)
	})
	if err != nil {
		return 0, errors.Wrapf(err, "[CountTotalProduct]query total product failed")
	}

	return count, nil
}

func (m *customProductModel) SFInsert(ctx context.Context, product *Product) (int64, error) {
	_, err := m.Insert(ctx, product)
	if err != nil {
		return 0, err
	}

	threading.GoSafe(func() {
		_ = util.Retry("del cache", 3, 3, func() error {
			return m.CachedConn.DelCache(CacheKeyProductCount(&product.StarId))
		})
	})

	return product.Id, nil
}

func (m *customProductModel) DeleteProduct(ctx context.Context, pId int64) error {
	err := m.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		_, err := s.ExecCtx(ctx, "delete from product_stock where product_id = ?", pId)
		if err != nil {
			return err
		}

		_, err = s.ExecCtx(ctx, "delete from product_color_detail where product_id = ?", pId)
		if err != nil {
			return err
		}

		_, err = s.ExecCtx(ctx, "delete from product where id = ?", pId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "[DeleteProduct] delete product failed")
	}

	threading.GoSafe(func() {
		_ = mr.Finish(func() error {
			_ = m.DelCache(CacheKeyProductCount(&pId))
			return nil
		}, func() error {
			_ = m.DelCache(CacheKeyProductColorList(pId))
			return nil
		}, func() error {
			_ = m.DelCache(CacheKeyProduct(pId))
			return nil
		})
	})

	return nil
}
