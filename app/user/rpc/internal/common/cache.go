package common

import (
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"ymir.com/pkg/vars"
	"ymir.com/pkg/xerr"
)

func CheckCapchaSendTimes(redis *redis.Redis, key string) error {
	sendTimesRaw, err := redis.Get(key)
	if err != nil {
		logx.Errorf("[SendCaptchaToEmail] get captcha send times failed, key: %s, err: %+v", key, err)
	}
	if len(sendTimesRaw) != 0 {
		sendTimes, _ := strconv.ParseInt(sendTimesRaw, 10, 64)
		if sendTimes+1 > vars.CaptchaMaxSendTimesPerDay {
			return xerr.NewErrCode(xerr.MaxCaptchaSendTimeError)
		}
		if _, err := redis.Incr(key); err != nil {
			logx.Errorf("[SendCaptchaToEmail] incr captcha send times failed, key: %s, err: %+v", key, err)
		}
	} else {
		if err := redis.Setex(key, "1", vars.CacheExpireIn1d); err != nil {
			logx.Errorf("[SendCaptchaToEmail] set captcha send times failed, key: %s, err: %+v", key, err)
		}
	}

	return nil
}
