package model

import (
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var ErrNotFound = sqlx.ErrNotFound

// Product Cache

func CacheKeyProductColorList(pId int64) string {
	return fmt.Sprintf("cache:productColorDetail:productId:%d", pId)
}

func CacheKeyProductCommentList(commentId int64) string {
	return fmt.Sprintf("cache:ymir:list:product:comment:%d", commentId)
}

func CacheKeyProductCount(starId *int64) string {
	if starId != nil {
		return fmt.Sprintf("cache:total:product:star:%d", *starId)
	}
	return "cache:total:product"
}

func CacheKeyProductBelongToStar(pId int64) string {
	return fmt.Sprintf("cache:starId:of:product:%d", pId)
}

func CacheKeyProduct(pId int64) string {
	return fmt.Sprintf("%s:%d",cacheYmirProductIdPrefix, pId)
}

func CacheKeyProductColor(pId int64) string {
	return fmt.Sprintf("%s:%d", cacheYmirProductColorDetailIdPrefix,pId)
}

func CacheKeyProductStock(pId int64) string {
	return fmt.Sprintf("%s:%d", cacheYmirProductStockIdPrefix, pId)
}
