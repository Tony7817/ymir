package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/threading"
)

var _ OrderModel = (*customOrderModel)(nil)

type (
	// OrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderModel.
	OrderModel interface {
		orderModel
		SFInsertTx(ctx context.Context, tx *sql.Tx, o *Order) (int64, error)
		DeleteTx(ctx context.Context, tx *sql.Tx, oId int64) error
		FindOrderList(ctx context.Context, userId int64, offset, limit int64) ([]*Order, error)
		CountOrderList(ctx context.Context, userId int64) (int64, error)
		UpdateOrderPartial(ctx context.Context, oId int64, uId int64, data *OrderPartial) error
	}

	customOrderModel struct {
		*defaultOrderModel
	}
)

type OrderPartial struct {
	Staus         *string         `db:"status"`
	TotalPrice    *int64          `db:"total_price"`
	Unit          *string         `db:"unit"`
	PaypalOrderId *sql.NullString `db:"paypal_order_id"`
	CancelReason  *sql.NullString `db:"cancel_reason"`
}

// NewOrderModel returns a model for the database table.
func NewOrderModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrderModel {
	return &customOrderModel{
		defaultOrderModel: newOrderModel(conn, c, opts...),
	}
}

func (m *customOrderModel) SFInsertTx(ctx context.Context, tx *sql.Tx, o *Order) (int64, error) {
	_, err := tx.ExecContext(ctx, "insert into `order` (id, user_id, request_id, status, total_price) values (?,?,?,?,?)", o.Id, o.UserId, o.RequestId, o.Status, o.TotalPrice)
	if err != nil {
		return 0, errors.Wrap(err, "[SFInsert] insert order failed")
	}

	threading.GoSafe(func() {
		_ = m.DelCache(CacheOrderListTotalCount(o.UserId))
	})

	return o.Id, nil
}

func (m *customOrderModel) FindOrderList(ctx context.Context, userId int64, offset, limit int64) ([]*Order, error) {
	var list []*Order
	err := m.QueryRowsNoCacheCtx(ctx, &list, "select * from `order` where user_id = ? order by created_at desc limit ?,?", userId, offset, limit)
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

func (m *customOrderModel) DeleteTx(ctx context.Context, tx *sql.Tx, oId int64) error {
	_, err := tx.ExecContext(ctx, "delete from `order` where id = ?", oId)
	if err != nil {
		return errors.Wrap(err, "[DeleteTx] delete order failed")
	}

	return nil
}

func (m *customOrderModel) UpdateOrderPartial(ctx context.Context, oId int64, uId int64, data *OrderPartial) error {
	var query = "update `order` set"
	var args []any
	if data.Staus != nil {
		query += " status = ?,"
		args = append(args, *data.Staus)
	}
	if data.TotalPrice != nil {
		query += " total_price = ?,"
		args = append(args, *data.TotalPrice)
	}
	if data.Unit != nil {
		query += " unit = ?,"
		args = append(args, *data.Unit)
	}
	if data.PaypalOrderId != nil {
		query += " paypal_order_id = ?,"
		args = append(args, *data.PaypalOrderId)
	}
	if data.CancelReason != nil {
		query += " cancel_reason = ?,"
		args = append(args, *data.CancelReason)
	}
	query = query[:len(query)-1]
	query += " where id = ?"
	args = append(args, oId)

	var cacheYmirOrderId = fmt.Sprintf("%s%v", cacheYmirOrderIdPrefix, oId)
	var cacheYmirOrderIdUserId = fmt.Sprintf("%s%v:%v", cacheYmirOrderIdUserIdPrefix, oId, uId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx, query, args...)
	}, cacheYmirOrderId, cacheYmirOrderIdUserId)
	if err != nil {
		return errors.Wrap(err, "[UpdateOrderPartial] update order failed")
	}

	return nil
}
