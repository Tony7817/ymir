package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest/httpx"
	"ymir.com/app/bffd/internal/config"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/pkg/result"
	"ymir.com/pkg/vars"
	"ymir.com/pkg/xerr"
)

type TimerMiddleware struct {
	Redis *redis.Redis
}

func NewTimerMiddleware(c config.Config) *TimerMiddleware {
	return &TimerMiddleware{
		Redis: redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
	}
}

func (m *TimerMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendSignupCaptchaRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		var cacheKey string
		if req.Email != nil {
			cacheKey = vars.GetCaptchaEmailLastRequestTimeKey(*req.Email)
		} else if req.Phonenumber != nil {
			cacheKey = vars.GetCaptchaPhonenumberLastRequestTimeKey(*req.Phonenumber)
		} else {
			result.HttpResult(r, w, nil, xerr.NewErrCode(xerr.ReuqestParamError))
		}

		lastReqTime, err := m.Redis.Get(cacheKey)
		if err != nil {
			logx.Errorf("[TimerMiddleware] get last request time failed, err: %+v", err)
			err = nil
		}
		if lastReqTime != "" {
			lastReqTimeInt, _ := strconv.ParseInt(lastReqTime, 10, 64)
			if time.Now().Unix()-lastReqTimeInt < int64(time.Minute) {
				result.HttpResult(r, w, nil, xerr.NewErrMsg("request too frequently"))
				return
			}
		} else {
			_, err = m.Redis.SetnxEx(cacheKey, "1", vars.CacheExpireIn1m)
			if err != nil {
				logx.Errorf("[TimerMiddleware] set last request time failed, err: %+v", err)
			}
		}

		var ctx = context.WithValue(r.Context(), vars.RequestContextKey, req)
		r = r.WithContext(ctx)
		next(w, r)
	}
}
