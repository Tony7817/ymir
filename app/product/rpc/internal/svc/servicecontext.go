package svc

import (
	"time"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ymir.com/app/product/rpc/internal/config"
	"ymir.com/app/product/rpc/model"
)

const localCacheExpire = time.Duration(time.Second * 60)

type ServiceContext struct {
	Config       config.Config
	ProductModel model.ProductModel
	LocalCache   *collection.Cache
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	localCache, err := collection.NewCache(localCacheExpire)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		ProductModel: model.NewProductModel(conn, c.CacheRedis),
		LocalCache:   localCache,
		Config:       c,
	}
}
