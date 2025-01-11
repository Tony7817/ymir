package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/threading"

	"ymir.com/pkg/id"
	"ymir.com/pkg/vars"
)

var _ ProductColorDetailModel = (*customProductColorDetailModel)(nil)

type (
	// ProductColorDetailModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductColorDetailModel.
	ProductColorDetailModel interface {
		productColorDetailModel
		FindColorListByPId(ctx context.Context, productId int64) ([]ProductColorDetail, error)
		SFInsert(ctx context.Context, color *ProductColorDetail) (int64, error)
		UpdatePartial(ctx context.Context, data *ProductColorDetailPartial) (sql.Result, error)
	}

	customProductColorDetailModel struct {
		*defaultProductColorDetailModel
		sf *id.Snowflake
	}
)

type ProductColorDetailPartial struct {
	Id            int64          `db:"id"`
	ProductId     *int64         `db:"product_id"`
	Color         *string        `db:"color"`
	Images        *string        `db:"images"`
	DetailImages  *string        `db:"detail_images"`
	Price         *float64       `db:"price"`
	Unit          *string        `db:"unit"`
	SoldNum       *sql.NullInt64 `db:"sold_num"`
	AvailableSize *string        `db:"available_size"`
	CoverUrl      *string        `db:"cover_url"`
}

// NewProductColorDetailModel returns a model for the database table.
func NewProductColorDetailModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ProductColorDetailModel {
	return &customProductColorDetailModel{
		defaultProductColorDetailModel: newProductColorDetailModel(conn, c, opts...),
		sf:                             id.NewSnowFlake(),
	}
}

func (m *customProductColorDetailModel) FindColorListByPId(ctx context.Context, productId int64) ([]ProductColorDetail, error) {
	var list []ProductColorDetail
	var cacheKey = CacheKeyProductColorList(productId)
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

func (m *customProductColorDetailModel) SFInsert(ctx context.Context, color *ProductColorDetail) (int64, error) {
	color.Id = m.sf.GenerateID()
	if _, err := m.Insert(ctx, color); err != nil {
		return 0, errors.Wrap(err, "[SFInsert] insert color error")
	}

	threading.GoSafe(func() {
		_ = m.DelCache(CacheKeyProductColorList(color.ProductId))
	})

	return color.Id, nil
}

func (m *customProductColorDetailModel) UpdatePartial(ctx context.Context, data *ProductColorDetailPartial) (sql.Result, error) {
	var query = "update product_color_detail set "
	var updateParams string
	var params []any
	if data.ProductId != nil {
		updateParams += "product_id = ?,"
		params = append(params, *data.ProductId)
	}
	if data.Color != nil {
		updateParams += "color = ?,"
		params = append(params, *data.Color)
	}
	if data.Images != nil {
		updateParams += "images = ?,"
		params = append(params, *data.Images)
	}
	if data.DetailImages != nil {
		updateParams += "detail_images = ?,"
		params = append(params, *data.DetailImages)
	}
	if data.SoldNum != nil {
		updateParams += "sold_num = ?,"
		params = append(params, *data.SoldNum)
	}
	if data.AvailableSize != nil {
		updateParams += "available_size = ?,"
		params = append(params, *data.AvailableSize)
	}
	if data.CoverUrl != nil {
		updateParams += "cover_url = ?,"
		params = append(params, *data.CoverUrl)
	}
	if updateParams == "" {
		return nil, errors.Wrapf(errors.New("update nothing is not allowded"), "[UpdatePartial] update fail")
	}

	updateParams = updateParams[:len(updateParams)-1]
	params = append(params, data.Id)

	var cond = " where id = ?"
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx, query+updateParams+cond, params...)
	}, fmt.Sprintf("%s%d", cacheYmirProductColorDetailIdPrefix, data.Id))
}
