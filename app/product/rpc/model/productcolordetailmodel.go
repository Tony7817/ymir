package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ymir.com/pkg/vars"
)

var _ ProductColorDetailModel = (*customProductColorDetailModel)(nil)

type (
	// ProductColorDetailModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductColorDetailModel.
	ProductColorDetailModel interface {
		productColorDetailModel
		FindColorListByPId(ctx context.Context, productId int64) ([]ProductColorDetail, error)
	}

	customProductColorDetailModel struct {
		*defaultProductColorDetailModel
	}
)

// NewProductColorDetailModel returns a model for the database table.
func NewProductColorDetailModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ProductColorDetailModel {
	return &customProductColorDetailModel{
		defaultProductColorDetailModel: newProductColorDetailModel(conn, c, opts...),
	}
}

func (m *customProductColorDetailModel) FindColorListByPId(ctx context.Context, productId int64) ([]ProductColorDetail, error) {
	var list []ProductColorDetail
	var cacheKey = fmt.Sprintf("cache:productColorDetail:productId:%d", productId)
	err := m.CachedConn.GetCacheCtx(ctx, cacheKey, &list)
	if err == nil {
		return list, nil
	} else if err != sql.ErrNoRows {
		logx.Errorf("[FindColorListByPId] get cache err: %v", err)
	}
	err = m.QueryRowsNoCacheCtx(ctx, &list, fmt.Sprintf("select * from %s where product_id = ?", m.table), productId)
	if err != nil {
		return nil, err
	}

	err = m.CachedConn.SetCacheWithExpireCtx(ctx, cacheKey, list, vars.CacheExpireIn1W)
	if err != nil {
		logx.Errorf("[FindColorListByPId] set cache err: %v", err)
	}

	return list, nil
}
