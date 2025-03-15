package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PaypalModel = (*customPaypalModel)(nil)

type (
	// PaypalModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPaypalModel.
	PaypalModel interface {
		paypalModel
		InsertDuplicateUpdate(ctx context.Context, data *Paypal) (sql.Result, error)
	}

	customPaypalModel struct {
		*defaultPaypalModel
	}
)

// NewPaypalModel returns a model for the database table.
func NewPaypalModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PaypalModel {
	return &customPaypalModel{
		defaultPaypalModel: newPaypalModel(conn, c, opts...),
	}
}

func (m *customPaypalModel) InsertDuplicateUpdate(ctx context.Context, data *Paypal) (sql.Result, error) {
	var cachePaypalId = fmt.Sprintf("%s%d", cacheYmirPaypalIdPrefix, data.Id)
	var cacheYmirPaypalOrderIdUserIdPrefix = fmt.Sprintf("%s%d:%d", cacheYmirPaypalOrderIdUserIdPrefix, data.OrderId, data.UserId)
	var cacheYmirPaypalRequestIdPrefix = fmt.Sprintf("%s%s", cacheYmirPaypalRequestIdPrefix, data.RequestId)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx, "insert into `paypal` (id, order_id, request_id, req_body, resp_body, paypal_order_id, user_id, err_msg) values (?,?,?,?,?,?,?,?) on duplicate key update req_body=?, resp_body=?, paypal_order_id=?, err_msg=?", data.Id, data.OrderId, data.RequestId, data.ReqBody, data.RespBody, data.PaypalOrderId, data.UserId, data.ErrMsg, data.ReqBody, data.RespBody, data.PaypalOrderId, data.ErrMsg)
	}, cachePaypalId, cacheYmirPaypalOrderIdUserIdPrefix, cacheYmirPaypalRequestIdPrefix)
}
