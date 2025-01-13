package model

import (
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var ErrNotFound = sqlx.ErrNotFound

func CacheKeyStarId(sId int64) string {
	return fmt.Sprintf("%s%d", cacheYmirStarIdPrefix, sId)
}
