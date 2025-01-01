package model

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ProductCommentImageModel = (*customProductCommentImageModel)(nil)

type (
	// ProductCommentImageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProductCommentImageModel.
	ProductCommentImageModel interface {
		productCommentImageModel
		FindImagesbyCommentId(ctx context.Context, commentId int64) ([]ProductCommentImage, error)
	}

	customProductCommentImageModel/*  */ struct {
		*defaultProductCommentImageModel
	}
)

// NewProductCommentImageModel returns a model for the database table.
func NewProductCommentImageModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ProductCommentImageModel {
	return &customProductCommentImageModel{
		defaultProductCommentImageModel: newProductCommentImageModel(conn, c, opts...),
	}
}

func (m *customProductCommentImageModel) FindImagesbyCommentId(ctx context.Context, commentId int64) ([]ProductCommentImage, error) {
	var res []ProductCommentImage
	err := m.QueryRowsNoCacheCtx(ctx, &res, "select * from product_comment_image where comment_id = ?", commentId)
	if err != nil {
		return nil, errors.Wrap(err, "[FindImagesbyCommentId] failed to get product comment image list")
	}

	return res, nil
}
