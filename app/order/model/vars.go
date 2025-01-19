package model

import (
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var ErrNotFound = sqlx.ErrNotFound
var (
	OrderStatusPending   = "pending"
	OrderStatusPaid      = "paid"
	OrderStatusCanceled  = "cancelled"
	OrderStatusShiped    = "shipped"
	OrderStatusCompleted = "completed"
)

func CacheOrderItemNotSoftDelete(id int64) string {
	return fmt.Sprintf("cache:ymir:orderItem:not:delete:id:%d", id)
}
