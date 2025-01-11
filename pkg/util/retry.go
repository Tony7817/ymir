package util

import (
	"time"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

func Retry(tag string, attempts int, sleep time.Duration, fn func() error) error {
	for i := 0; i < attempts; i++ {
		if err := fn(); err == nil {
			return nil
		}
		time.Sleep(sleep)
		logx.Errorf("[%s]Attempt %d failed, retrying...", tag, i+1)
	}
	logx.Errorf("[%s]retry fail %d times", tag, attempts)
	return errors.Wrapf(errors.New("retry fail by "+tag), "")
}
