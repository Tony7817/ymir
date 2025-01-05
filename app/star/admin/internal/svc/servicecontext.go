package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ymir.com/app/star/admin/internal/config"
	"ymir.com/app/star/model"
)

type ServiceContext struct {
	Config    config.Config
	StarModel model.StarModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		StarModel: model.NewStarModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
	}
}
