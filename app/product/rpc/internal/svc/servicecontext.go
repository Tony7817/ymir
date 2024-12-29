package svc

import (
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ymir.com/app/product/rpc/internal/config"
	"ymir.com/app/product/rpc/model"
)

const localCacheExpire = time.Duration(time.Second * 60)

type ServiceContext struct {
	Config                   config.Config
	ProductModel             model.ProductModel
	ProductCartModel         model.ProductCartModel
	ProductColorModel        model.ProductColorDetailModel
	ProductStockModel        model.ProductStockModel
	ProductCommentModel      model.ProductCommentModel
	ProductCommentImageModel model.ProductCommentImageModel
	Redis                    *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		ProductModel:             model.NewProductModel(conn, c.CacheRedis),
		ProductCartModel:         model.NewProductCartModel(conn, c.CacheRedis),
		ProductStockModel:        model.NewProductStockModel(conn, c.CacheRedis),
		ProductColorModel:        model.NewProductColorDetailModel(conn, c.CacheRedis),
		ProductCommentModel:      model.NewProductCommentModel(conn, c.CacheRedis),
		ProductCommentImageModel: model.NewProductCommentImageModel(conn, c.CacheRedis),
		Redis:                    redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
		Config:                   c,
	}
}
