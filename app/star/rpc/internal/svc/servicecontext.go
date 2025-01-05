package svc

import (
	"time"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ymir.com/app/star/model"
	"ymir.com/app/star/rpc/internal/config"
)

const localCacheExpire = time.Duration(time.Second * 60)

type ServiceContext struct {
	Config     config.Config
	StarModel  model.StarModel
	LocalCache *collection.Cache
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	localCache, err := collection.NewCache(localCacheExpire)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:     c,
		StarModel:  model.NewStarModel(conn, c.CacheRedis),
		LocalCache: localCache,
	}
}
