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
		SoftDeleteRollbackTx(ctx context.Context, tx *sql.Tx, id, userId int64) error
		FindOrderList(ctx context.Context, userId int64, offset, limit int64) ([]*Order, error)
		CountOrderList(ctx context.Context, userId int64) (int64, error)
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

	_, err = tx.ExecContext(ctx, "update order_item set is_delete = 1 where order_id = ?", id)
	if err != nil {
		return errors.Wrap(err, "[SoftDeleteTx] soft delete order item failed")
	}

	return nil
}

func (m *customOrderModel) SoftDeleteRollbackTx(ctx context.Context, tx *sql.Tx, id int64, userId int64) error {
	_, err := tx.ExecContext(ctx, "update `order` set is_delete = 0 where id = ? and user_id = ?", id)
	if err != nil {
		return errors.Wrap(err, "[SoftDeleteRollbackTx] soft delete rollback order failed")
	}

	_, err = tx.ExecContext(ctx, "update order_item set is_delete = 0 where order_id = ?", id)
	if err != nil {
		return errors.Wrap(err, "[SoftDeleteRollbackTx] soft delete rollback order item failed")
	}

	return nil
}

func (m *customOrderModel) FindOrderList(ctx context.Context, userId int64, offset, limit int64) ([]*Order, error) {
	var list []*Order
	err := m.QueryRowsNoCacheCtx(ctx, &list, "select * from `order` where user_id = ? limit ?,?", userId, offset, limit)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return list, nil
}

func (m *customOrderModel) CountOrderList(ctx context.Context, userId int64) (int64, error) {
	var count int64
	err := m.QueryRowCtx(ctx, &count, CacheOrderListTotalCount(userId), func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		return conn.QueryRowCtx(ctx, &count, "select count(*) from `order` where user_id = ?", userId)
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}
