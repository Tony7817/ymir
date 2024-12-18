package model

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ymir.com/pkg/vars"
)

var _ CaptchaModel = (*customCaptchaModel)(nil)

type (
	// CaptchaModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCaptchaModel.
	CaptchaModel interface {
		captchaModel
		FindCaptchaByEmail(email string) (string, error)
		FindCaptchaByPhonenumber(phonenumber string) (string, error)
	}

	customCaptchaModel struct {
		*defaultCaptchaModel
	}
)

// NewCaptchaModel returns a model for the database table.
func NewCaptchaModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CaptchaModel {
	return &customCaptchaModel{
		defaultCaptchaModel: newCaptchaModel(conn, c, opts...),
	}
}

func (m *customCaptchaModel) FindCaptchaByEmail(email string) (string, error) {
	var cacheKey = vars.GetEmailCapchaCacheKey(email)
	var captcha Captcha
	err := m.QueryRow(&captcha, cacheKey, func(conn sqlx.SqlConn, v any) error {
		return conn.QueryRow(&captcha, "select verify_code from captcha where email = ? and is_delete = 0 and created_at > DATE_SUB(now(), INTERVAL 5 MINUTE) order by created_at desc limit 1", email)
	})
	if err != nil {
		logx.Errorf("[FindCaptchaByEmail] find captcha by email failed, err: %+v", err)
		return "", err
	}

	return captcha.VerifyCode, nil
}

func (m *customCaptchaModel) FindCaptchaByPhonenumber(phonenumber string) (string, error) {
	var cacheKey = vars.GetPhonenumberCapchaCacheKey(phonenumber)
	var captcha Captcha
	err := m.QueryRow(&captcha, cacheKey, func(conn sqlx.SqlConn, v any) error {
		return conn.QueryRow(&captcha, "select verify_code from captcha where phone_number = ? and is_delete = 0 and created_at > DATE_SUB(now(), INTERVAL 5 MINUTE) order by created_at desc limit 1", phonenumber)
	})
	if err != nil {
		logx.Errorf("[FindCaptchaByEmail] find captcha by email failed, err: %+v", err)
		return "", err
	}

	return captcha.VerifyCode, nil
}
