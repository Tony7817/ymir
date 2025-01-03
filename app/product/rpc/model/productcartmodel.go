package model

import (
	"context"
	"database/sql"
	"fmt"

	"ymir.com/pkg/id"
	"ymir.com/pkg/xerr"

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
		IncrProductAmount(ctx context.Context, userId int64, productId int64, color string) error
		DescProductAmount(ctx context.Context, userId int64, productId int64, color string) error
		AddToProductCart(ctx context.Context, userId int64, productId int64, size string, colorId int64) (sql.Result, error)
		DeleteProductFromCart(ctx context.Context, userId int64, productId int64, colorId int64) error
	}

	customProductCartModel struct {
		*defaultProductCartModel
	}
)

const (
	// %s: userId
	cacheYmirProductCartTotalUserId = "cache:ymir:productCart:total:%d"
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
	var cacheKey = fmt.Sprintf(cacheYmirProductCartTotalUserId, userId)
	err := m.QueryRowCtx(ctx, &count, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		return conn.QueryRow(&count, "select count(*) from product_cart where user_id = ?", userId)
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *customProductCartModel) IncrProductAmount(ctx context.Context, userId int64, productId int64, color string) error {
	err := m.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		var inStock int64
		err := s.QueryRowCtx(ctx, &inStock, "select in_stock from product_stock where product_id = ? and color = ?", productId, color)
		if err != nil {
			return errors.Wrapf(err, "[IncrProductAmount] failed to get product stock")
		}

		if inStock <= 0 {
			return xerr.NewErrCode(xerr.CaptchaExpireError)
		}

		_, err = s.ExecCtx(ctx, "update product_cart set amount = amount + 1 where user_id = ? and product_id = ? and color = ?", userId, productId, color)
		if err != nil {
			return errors.Wrapf(err, "[IncrProductAmount] failed to increase product amount")
		}

		return nil
	})
	if err != nil {
		return err
	}

	// delete cachhe
	threading.GoSafe(func() {
		_ = m.CachedConn.DelCache(fmt.Sprintf(cacheYmirProductCartTotalUserId, userId))
	})

	return nil
}

func (m *customProductCartModel) DescProductAmount(ctx context.Context, userId int64, productId int64, color string) error {
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx, "update product_cart set amount = amount - 1 where user_id = ? and product_id = ? and color = ?", userId, productId, color)
	}, fmt.Sprintf(cacheYmirProductCartTotalUserId, userId))
	return err
}

func (m *customProductCartModel) AddToProductCart(ctx context.Context, userId int64, productId int64, size string, colorId int64) (sql.Result, error) {
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx, "insert into product_cart (id, user_id, product_id, size, amount, color_id) values (?, ?, ?, 1, ?) on duplicate key update amount = amount + 1", id.SF.GenerateID(), userId, productId, size, colorId)
	}, fmt.Sprintf(cacheYmirProductCartTotalUserId, userId))
}

func (m *customProductCartModel) DeleteProductFromCart(ctx context.Context, productId int64, colorId int64, userId int64) error {
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx, "delete from product_cart where product_id = ? and color_id = ? and user_id = ?", productId, colorId, userId)
	}, fmt.Sprintf(cacheYmirProductCartTotalUserId, userId))
	if err != nil {
		return err
	}

	return nil
}
