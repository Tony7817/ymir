package model

import (
	"context"
	"database/sql"

	"ymir.com/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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
		AddToProductCart(ctx context.Context, userId int64, productId int64, size string, color int64) (sql.Result, error)
		SoftRemoveToProductFromCart(ctx context.Context, userId int64, productId int64, color string) (sql.Result, error)
	}

	customProductCartModel struct {
		*defaultProductCartModel
	}
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
		return nil, err
	}

	return res, nil
}

func (m *customProductCartModel) CountTotalProductOfUser(ctx context.Context, userId int64) (int64, error) {
	var count int64
	var cacheKey = "cache:productCart:userId:total"
	err := m.QueryRowCtx(ctx, &count, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		return conn.QueryRow(ctx, "select count(*) from product_cart where user_id = ?", userId)
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

	return nil
}

func (m *customProductCartModel) DescProductAmount(ctx context.Context, userId int64, productId int64, color string) error {
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx, "update product_cart set amount = amount - 1 where user_id = ? and product_id = ? and color = ?", userId, productId, color)
	})
	return err
}

func (m *customProductCartModel) AddToProductCart(ctx context.Context, userId int64, productId int64, size string, colorId int64) (sql.Result, error) {
	return m.ExecNoCacheCtx(ctx, "insert into product_cart (user_id, product_id, size, amount, color_id) values (?, ?, ?, 1, ?) on duplicate key update amount = amount + 1", userId, productId, size, colorId)
}

func (m *customProductCartModel) SoftRemoveToProductFromCart(ctx context.Context, userId int64, productId int64, color string) (sql.Result, error) {
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx, "update product_cart set is_delete = 1 where user_id = ? and product_id = ? and color = ?", userId, productId, color)
	})
}
