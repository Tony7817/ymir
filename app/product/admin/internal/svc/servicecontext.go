package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ymir.com/app/product/admin/internal/config"
	"ymir.com/app/product/model"
)

type ServiceContext struct {
	Config            config.Config
	ProductModel      model.ProductModel
	ProductColorModel model.ProductColorDetailModel
	ProductStockModel model.ProductStockModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:            c,
		ProductModel:      model.NewProductModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		ProductColorModel: model.NewProductColorDetailModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		ProductStockModel: model.NewProductStockModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
	}
}
