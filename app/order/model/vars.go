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

func CacheOrderListTotalCount(userId int64) string {
	return fmt.Sprintf("cache:ymir:order:list:total:count:uid:%d", userId)
}

func CacheOrderItemsByOrderId(orderId int64) string {
	return fmt.Sprintf("cache:ymir:order:items:oid:%d", orderId)
}

func CacheOrderItemByOrderIdUserId(orderId, userId int64) string {
	return fmt.Sprintf("cache:ymir:order:item:oid:%d:uid:%d", orderId, userId)
}

func CacheCountOrderItemByOrderIdUserId(orderId, userId int64) string {
	return fmt.Sprintf("cache:ymir:order:item:count:oid:%d:uid:%d", orderId, userId)
}