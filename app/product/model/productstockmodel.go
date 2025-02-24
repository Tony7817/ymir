package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ymir.com/pkg/xerr"
)

var _ ProductStockModel = (*customProductStockModel)(nil)

type (
	// ProductStockModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductStockModel.
	ProductStockModel interface {
		productStockModel
		InsertSF(ctx context.Context, data *ProductStock) (int64, error)
		DecreaseProductStockTx(ctx context.Context, tx *sql.Tx, pId, cId, quantity int64, size string) error
		IncreaseProductStockTx(ctx context.Context, tx *sql.Tx, pId, cId, quantity int64, size string) error
	}

	customProductStockModel struct {
		*defaultProductStockModel
	}
)

// NewProductStockModel returns a model for the database table.
func NewProductStockModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ProductStockModel {
	return &customProductStockModel{
		defaultProductStockModel: newProductStockModel(conn, c, opts...),
	}
}

func (m *customProductStockModel) InsertSF(ctx context.Context, data *ProductStock) (int64, error) {
	if _, err := m.Insert(ctx, data); err != nil {
		return 0, errors.Wrapf(err, "[InsertSF]insert product stock fail")
	}

	return data.Id, nil
}

func (m *customProductStockModel) DecreaseProductStockTx(ctx context.Context, tx *sql.Tx, pId, cId, quantity int64, size string) error {
	ret, err := tx.ExecContext(ctx, "update product_stock set in_stock = in_stock - ? where product_id = ? and color_id = ? and size = ? and in_stock >= ?", quantity, pId, cId, size, quantity)
	if err != nil {
		return errors.Wrap(err, "[DecreaseProductStockTx]update product stock fail")
	}

	affected, err := ret.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "[DecreaseProductStockTx]get affected rows fail")
	}
	if affected == 0 {
		return errors.Wrap(xerr.NewErrCode(xerr.DbUpdateAffectedZeroError), "[DecreaseProductStockTx]stock not enough")
	}

	_ = m.DelCache(fmt.Sprintf("%s%v:%v:%v", cacheYmirProductStockProductIdColorIdSizePrefix, pId, cId, size))

	return nil
}

func (m *customProductStockModel) IncreaseProductStockTx(ctx context.Context, tx *sql.Tx, pId, cId, quantity int64, size string) error {
	_, err := tx.ExecContext(ctx, "update product_stock set in_stock = in_stock + ? where product_id = ? and color_id = ? and size = ?", quantity, pId, cId, size)
	if err != nil {
		return errors.Wrap(err, "[IncreaseProductStockTx]update product stock fail")
	}

	_ = m.DelCache(fmt.Sprintf("%s%v:%v:%v", cacheYmirProductStockProductIdColorIdSizePrefix, pId, cId, size))

	return nil
}
