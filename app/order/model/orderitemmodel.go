package model

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/threading"
)

var _ OrderItemModel = (*customOrderItemModel)(nil)

type (
	// OrderItemModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOrderItemModel.
	OrderItemModel interface {
		orderItemModel
		FindOneByOrderId(ctx context.Context, id int64) ([]OrderItem, error)
		FindOrderItemsByOrderIdUserId(ctx context.Context, orderId int64, userId int64) ([]OrderItem, error)
		CountOrderItemsByOrderIdUserId(ctx context.Context, orderId int64, userId int64) (int64, error)
		SFInsertTx(ctx context.Context, tx *sql.Tx, oitem *OrderItem) (int64, error)
		DeleteTx(ctx context.Context, tx *sql.Tx, oId int64) error
	}

	customOrderItemModel struct {
		*defaultOrderItemModel
	}
)

// NewOrderItemModel returns a model for the database table.
func NewOrderItemModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrderItemModel {
	return &customOrderItemModel{
		defaultOrderItemModel: newOrderItemModel(conn, c, opts...),
	}
}

func (m *customOrderItemModel) SFInsertTx(ctx context.Context, tx *sql.Tx, oitem *OrderItem) (int64, error) {
	_, err := tx.ExecContext(ctx, "insert into order_item (id, user_id, order_id, product_id, product_description, product_color_id, product_color, product_color_cover_url, size, quantity, price, subtotal) values (?,?,?,?,?,?,?,?,?,?,?,?)",
		oitem.Id, oitem.UserId, oitem.OrderId, oitem.ProductId, oitem.ProductDescription, oitem.ProductColorId, oitem.ProductColor, oitem.ProductColorCoverUrl, oitem.Size, oitem.Quantity, oitem.Price, oitem.Subtotal)
	if err != nil {
		return 0, errors.Wrap(err, "[SFInsert] insert order item fail")
	}
	return 0, nil
}

func (m *customOrderItemModel) FindOneByOrderId(ctx context.Context, orderId int64) ([]OrderItem, error) {
	var o []OrderItem
	_ = m.CachedConn.GetCacheCtx(ctx, CacheOrderItemsByOrderId(orderId), &o)
	if len(o) > 0 {
		return o, nil
	}

	err := m.QueryRowsNoCacheCtx(ctx, &o, "select * from order_item where order_id = ?", orderId)
	if err != nil {
		return nil, err
	}

	threading.GoSafe(func() {
		_ = m.CachedConn.SetCache(CacheOrderItemsByOrderId(orderId), o)
	})

	return o, nil
}

func (m *customOrderItemModel) DeleteTx(ctx context.Context, tx *sql.Tx, oId int64) error {
	_, err := tx.ExecContext(ctx, "delete from order_item where id = ?", oId)
	if err != nil {
		return errors.Wrap(err, "[DeleteOrderItemTx] delete order_item failed")
	}

	return nil
}

func (m *customOrderItemModel) FindOrderItemsByOrderIdUserId(ctx context.Context, orderId int64, userId int64) ([]OrderItem, error) {
	var o []OrderItem

	_ = m.CachedConn.GetCacheCtx(ctx, CacheOrderItemByOrderIdUserId(orderId, userId), &o)
	if len(o) > 0 {
		return o, nil
	}

	err := m.QueryRowsNoCacheCtx(ctx, &o, "select * from order_item where order_id = ? and user_id = ?", orderId, userId)
	if err != nil {
		return nil, err
	}

	threading.GoSafe(func() {
		_ = m.CachedConn.SetCache(CacheOrderItemByOrderIdUserId(orderId, userId), o)
	})

	return o, nil
}

func (m *customOrderItemModel) CountOrderItemsByOrderIdUserId(ctx context.Context, orderId int64, userId int64) (int64, error) {
	var count int64
	err := m.QueryRowCtx(ctx, &count, CacheCountOrderItemByOrderIdUserId(orderId, userId), func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		return conn.QueryRowCtx(ctx, &count, "select count(*) from order_item where order_id = ? and user_id = ?", orderId, userId)
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}
