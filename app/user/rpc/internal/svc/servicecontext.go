package svc

import (
	"ymir.com/app/user/rpc/internal/config"
	"ymir.com/app/user/rpc/model"
	"ymir.com/pkg/aliyun"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	UserModel       model.UserModel
	CaptchaModel    model.CaptchaModel
	UserGoogleModel model.UserGoogleModel
	UserLocalModel  model.UserLocalModel

	Redis *redis.Redis

	// aliyun
	EmailClient *aliyun.ClientWrapper
}

func NewServiceContext(c config.Config) *ServiceContext {
	client, err := aliyun.NewClientWrapper()
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:          c,
		UserModel:       model.NewUserModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		CaptchaModel:    model.NewCaptchaModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		UserGoogleModel: model.NewUserGoogleModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		UserLocalModel:  model.NewUserLocalModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		Redis:           redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
		EmailClient:     client,
	}
}
