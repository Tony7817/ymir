package model

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ProductCommentModel = (*customProductCommentModel)(nil)
var CacheYmirProductCommentTotalCount = "cache:ymir:productComment:total:%d"

type (
	// ProductCommentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductCommentModel.
	ProductCommentModel interface {
		productCommentModel
		CountTotalProductComment(ctx context.Context, productId int64) (int64, error)
		FindProductCommentList(ctx context.Context, productId int64, page int64, pageSize int64) ([]ProductComment, error)
	}

	customProductCommentModel struct {
		*defaultProductCommentModel
	}
)

// NewProductCommentModel returns a model for the database table.
func NewProductCommentModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ProductCommentModel {
	return &customProductCommentModel{
		defaultProductCommentModel: newProductCommentModel(conn, c, opts...),
	}
}

func (m *customProductCommentModel) CountTotalProductComment(ctx context.Context, productId int64) (int64, error) {
	var count int64
	err := m.QueryRowCtx(ctx, &count, fmt.Sprintf(CacheYmirProductCommentTotalCount, productId), func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		return conn.QueryRowCtx(ctx, &count, "select count(*) from product_comment where product_id = ?", productId)
	})
	if err != nil {
		return 0, errors.Wrapf(err, "[CountTotalProductComment] failed to get total count")
	}

	return count, nil
}

func (m *customProductCommentModel) FindProductCommentList(ctx context.Context, productId int64, page int64, pageSize int64) ([]ProductComment, error) {
	var (
		res    []ProductComment
		offset = (page - 1) * pageSize
	)
	err := m.QueryRowsNoCache(&res, "select * from product_comment where product_id = ? limit ?,?", productId, offset, pageSize)
	if err != nil {
		return nil, errors.Wrap(err, "[FindProductCommentList] failed to get product comment list")
	}

	return res, err
}
