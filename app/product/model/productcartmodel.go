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

var _ ProductCartModel = (*customProductCartModel)(nil)

type (
	// ProductCartModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductCartModel.
	ProductCartModel interface {
		productCartModel
		FindProductsOfUser(userId int64, page int64, pageSize int64) ([]*ProductCart, error)
		CountTotalProductOfUser(ctx context.Context, userId int64) (int64, error)
		IncrProductAmount(ctx context.Context, userId int64, productId int64, colorId int64, size string) error
		DescProductAmount(ctx context.Context, userId int64, productId int64, colorId int64, size string) error
		InsertIntoProductCart(ctx context.Context, pcId int64, userId int64, productId int64, size string, colorId int64) error
		DeleteProductFromCart(ctx context.Context, userId int64, productId int64, colorId int64) error
		CheckoutProduct(ctx context.Context, tx *sql.Tx, oId int64, id int64, pc *ProductCart) error
		CheckoutProductRollback(ctx context.Context, tx *sql.Tx, pc *ProductCartCheckout) error
	}

	customProductCartModel struct {
		*defaultProductCartModel
	}
)

const (
// %s: userId
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
	if err := m.QueryRowsNoCache(&res, "select * from product_cart where user_id = ? limit ?, ?", userId, offset, pageSize); err != nil {
		return nil, errors.Wrapf(err, "[FindProductsOfUser] failed to get product cart list")
	}

	return res, nil
}

func (m *customProductCartModel) CountTotalProductOfUser(ctx context.Context, userId int64) (int64, error) {
	var count int64
	var cacheKey = CacheYmirProductCartTotalUserId(userId)
	err := m.QueryRowCtx(ctx, &count, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		return conn.QueryRow(&count, "select count(*) from product_cart where user_id = ?", userId)
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *customProductCartModel) IncrProductAmount(ctx context.Context, userId int64, productId int64, colorId int64, size string) error {
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx, "update product_cart set amount = amount + 1 where user_id = ? and product_id = ? and color_id = ? and size = ?", userId, productId, colorId, size)
	}, fmt.Sprintf("%s%v:%v:%v:%v", cacheYmirProductCartUserIdColorIdProductIdSizePrefix, userId, colorId, productId, size))
	return err
}

func (m *customProductCartModel) DescProductAmount(ctx context.Context, userId int64, productId int64, colorId int64, size string) error {
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx, "update product_cart set amount = amount - 1 where user_id = ? and color_id = ? and product_id = ? and size = ?", userId, colorId, productId, size)
	}, fmt.Sprintf("%s%v:%v:%v:%v", cacheYmirProductCartUserIdColorIdProductIdSizePrefix, userId, colorId, productId, size))
	return err
}

func (m *customProductCartModel) InsertIntoProductCart(ctx context.Context, pcId int64, userId int64, productId int64, size string, colorId int64) error {
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx, "insert into product_cart (id, user_id, product_id, size, amount, color_id) values (?, ?, ?, ?, 1, ?) on duplicate key update amount = amount + 1", pcId, userId, productId, size, colorId)
	}, CacheYmirProductCartTotalUserId(userId))
	if err != nil {
		return err
	}
	return nil
}

func (m *customProductCartModel) DeleteProductFromCart(ctx context.Context, productId int64, colorId int64, userId int64) error {
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx, "delete from product_cart where product_id = ? and color_id = ? and user_id = ?", productId, colorId, userId)
	}, CacheYmirProductCartTotalUserId(userId))
	if err != nil {
		return err
	}

	return nil
}

func (m *customProductCartModel) CheckoutProduct(ctx context.Context, tx *sql.Tx, oId int64, id int64, pc *ProductCart) error {
	_, err := tx.ExecContext(ctx, "insert into product_cart_checkout (id, product_cart_id, origin_created_at, order_id, user_id, product_id, color_id, size, amount) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", id, pc.Id, pc.CreatedAt, oId, pc.UserId, pc.ProductId, pc.ColorId, pc.Size, pc.Amount)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "delete from product_cart where id = ?", pc.Id)
	if err != nil {
		return err
	}

	threading.GoSafe(func() {
		m.DelCache(fmt.Sprintf("%s%v:%v:%v:%v", cacheYmirProductCartUserIdColorIdProductIdSizePrefix, pc.UserId, pc.ColorId, pc.ProductId, pc.Size))
		m.DelCache(fmt.Sprintf("%s%v", cacheYmirProductCartIdPrefix, id))
		m.DelCache(CacheYmirProductCartTotalUserId(pc.UserId))
	})

	return nil
}

func (m *customProductCartModel) CheckoutProductRollback(ctx context.Context, tx *sql.Tx, pc *ProductCartCheckout) error {
	_, err := tx.ExecContext(ctx, "insert into product_cart (id, created_at, updated_at, user_id, product_id, color_id, size, amount) values (?, ?, ?, ?, ?, ?, ?, ?)", pc.Id, pc.OriginCreatedAt, pc.UpdatedAt, pc.UserId, pc.ProductId, pc.ColorId, pc.Size, pc.Amount)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "delete from product_cart_checkout where id = ?", pc.Id)
	if err != nil {
		return err
	}

	return nil
}
