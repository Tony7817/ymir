package model

import (
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var ErrNotFound = sqlx.ErrNotFound

func CacheKeyProductColorList(pId int64) string {
	return fmt.Sprintf("cache:productColorDetail:productId:%d", pId)
}

func CacheKeyProductCommentList(commentId int64) string {
	return fmt.Sprintf("cache:ymir:list:product:comment:%d", commentId)
}

func ProductCountCacheKey(starId *int64) string {
	if starId != nil {
		return fmt.Sprintf("cache:total:product:star:%d", *starId)
	}
	return "cache:total:product"
}
