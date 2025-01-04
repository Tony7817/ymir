package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ymir.com/app/user/admin/internal/config"
	"ymir.com/app/user/model"
)

type ServiceContext struct {
	Config         config.Config
	OrganizerModel model.OrganizerModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		OrganizerModel: model.NewOrganizerModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
	}
}
