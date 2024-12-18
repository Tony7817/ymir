package svc

import (
	"ymir.com/app/user/rpc/internal/config"
	"ymir.com/app/user/rpc/model"
	"ymir.com/pkg/aliyun"

	dm "github.com/alibabacloud-go/dm-20151123/v2/client"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config       config.Config
	UserModel    model.UserModel
	CaptchaModel model.CaptchaModel
	Redis        *redis.Redis

	// aliyun
	EmailClient *dm.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	client, err := aliyun.NewEmailClient()
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:       c,
		UserModel:    model.NewUserModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		CaptchaModel: model.NewCaptchaModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		Redis:        redis.New(c.Redis.Host, redis.WithPass(c.Redis.Pass)),
		EmailClient:  client,
	}
}
