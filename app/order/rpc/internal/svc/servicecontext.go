package svc

import (
	"database/sql"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ymir.com/app/order/model"
	"ymir.com/app/order/rpc/internal/config"
)

type ServiceContext struct {
	Config              config.Config
	OrderModel          model.OrderModel
	OrderItemModel      model.OrderItemModel
	OrderStatusLogModel model.OrderStatusLogModel
	PaypalModel         model.PaypalModel
	DB                  *sql.DB
	Redis               *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := sqlx.NewMysql(c.DataSource).RawDB()
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:              c,
		OrderModel:          model.NewOrderModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		OrderItemModel:      model.NewOrderItemModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		OrderStatusLogModel: model.NewOrderStatusLogModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		PaypalModel:         model.NewPaypalModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		DB:                  db,
		Redis:               redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
	}
}
