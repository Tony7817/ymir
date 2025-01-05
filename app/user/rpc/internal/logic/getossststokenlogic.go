package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"ymir.com/app/user/rpc/internal/svc"
	"ymir.com/app/user/rpc/user"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetOssStsTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOssStsTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOssStsTokenLogic {
	return &GetOssStsTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOssStsTokenLogic) GetOssStsToken(in *user.GetOssStsTokenRequest) (*user.GetOssStsTokenResponse, error) {
	var cacheKey = fmt.Sprintf("cache:oss:sts:token:%d", in.UserId)

	var res user.GetOssStsTokenResponse
	cache, err := l.svcCtx.Redis.GetCtx(l.ctx, cacheKey)
	if err != nil {
		l.Logger.Errorf("[GetOssStsToken] get cache fail, err:%+v", err)
	}
	if cache != "" {
		if err := json.Unmarshal([]byte(cache), &res); err == nil {
			return &res, nil
		}
		l.Logger.Errorf("[GetOssStsToken] unmarshal cache fail, err:%+v", err)
	}

	token, err := l.svcCtx.AliyunOssClient.GetSTSToken(fmt.Sprintf("%d", in.UserId))
	if err != nil {
		return nil, errors.Wrapf(err, "[GetOssSTSToken] get token fail, err:%+v", err)
	}
	l.Logger.Infof("token:%+v", *token)

	// token is already a pointer
	if err := copier.Copy(&res, token); err != nil {
		return nil, errors.Wrapf(err, "[GetOssSTSToken] copy fail, err:%+v", err)
	}

	// expire in 60 minites
	resRaw, err := json.Marshal(&res)
	if err != nil {
		logx.Errorf("[GetOssSTSToken] marshal fail, err:%+v", err)
	}
	if err := l.svcCtx.Redis.SetexCtx(l.ctx, cacheKey, string(resRaw), 60*59); err != nil {
		logx.Errorf("[GetOssSTSToken] set cache fail, err:%+v", err)
	}

	return &res, nil
}
