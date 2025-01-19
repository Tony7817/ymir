package model

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderModel = (*customOrderModel)(nil)

type (
	// OrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderModel.
	OrderModel interface {
		orderModel
		SFInsertTx(ctx context.Context, tx *sql.Tx, o *Order) (int64, error)
		SoftDeleteTx(ctx context.Context, tx *sql.Tx, id, userId int64) error
		DeleteTx(ctx context.Context, tx *sql.Tx, id, userId int64) error
	}

	customOrderModel struct {
		*defaultOrderModel
	}
)

// NewOrderModel returns a model for the database table.
func NewOrderModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrderModel {
	return &customOrderModel{
		defaultOrderModel: newOrderModel(conn, c, opts...),
	}
}

func (m *customOrderModel) SFInsertTx(ctx context.Context, tx *sql.Tx, o *Order) (int64, error) {
	_, err := tx.ExecContext(ctx, "insert into `order` (id, user_id, status, total_price) values (?,?,?,?)", o.Id, o.UserId, o.Status, o.TotalPrice)
	if err != nil {
		return 0, errors.Wrap(err, "[SFInsert] insert order failed")
	}

	return o.Id, nil
}

func (m *customOrderModel) SoftDeleteTx(ctx context.Context, tx *sql.Tx, id int64, userId int64) error {
	_, err := tx.ExecContext(ctx, "update `order` set is_delete = 1 where id = ? and user_id = ?", id)
	if err != nil {
		return errors.Wrap(err, "[SoftDeleteTx] soft delete order failed")
	}

	return nil
}

func (m *customOrderModel) DeleteTx(ctx context.Context, tx *sql.Tx, id int64, userId int64) error {
	_, err := tx.ExecContext(ctx, "delete from `order` where id = ? and user_id = ?", id, userId)
	if err != nil {
		return errors.Wrap(err, "[DeleteTx] delete order failed")
	}

	return nil
}
