package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ymir.com/app/product/admin/internal/config"
	"ymir.com/app/product/model"
)

type ServiceContext struct {
	Config       config.Config
	ProductModel model.ProductModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		ProductModel: model.NewProductModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
	}
}
