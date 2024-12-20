package model

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
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
		FindCaptchaByEmail(email string) (*Captcha, error)
		FindCaptchaByPhonenumber(phonenumber string) (*Captcha, error)
		InsertCaptchaToDbAndCache(captcha Captcha) (sql.Result, error)
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

func (m *customCaptchaModel) FindCaptchaByEmail(email string) (*Captcha, error) {
	var cacheKey = vars.GetEmailCapchaCacheKey(email)
	var captcha Captcha
	err := m.QueryRow(&captcha, cacheKey, func(conn sqlx.SqlConn, v any) error {
		return conn.QueryRow(&captcha, "select * from captcha where email = ? and is_delete = 0 and created_at > DATE_SUB(now(), INTERVAL 5 MINUTE) order by created_at desc limit 1", email)
	})
	if err != nil {
		logx.Errorf("[FindCaptchaByEmail] find captcha by email failed, err: %+v", err)
		return nil, errors.Wrapf(err, "[FindValidCaptchaByEmail] find captcha by email failed, email:%s", email)
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &captcha, nil
}

func (m *customCaptchaModel) FindCaptchaByPhonenumber(phonenumber string) (*Captcha, error) {
	var cacheKey = vars.GetPhonenumberCapchaCacheKey(phonenumber)
	var captcha Captcha
	err := m.QueryRow(&captcha, cacheKey, func(conn sqlx.SqlConn, v any) error {
		return conn.QueryRow(&captcha, "select * from captcha where phone_number = ? and is_delete = 0 and created_at > DATE_SUB(now(), INTERVAL 5 MINUTE) order by created_at desc limit 1", phonenumber)
	})
	if err != nil {
		logx.Errorf("[FindCaptchaByEmail] find captcha by email failed, err: %+v", err)
		return nil, errors.Wrapf(err, "[FindValidCaptchaByEmail] find captcha by email failed, phonenumber:%s", phonenumber)
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &captcha, nil
}

func (m *customCaptchaModel) InsertCaptchaToDbAndCache(captcha Captcha) (sql.Result, error) {
	var cacheKey string
	if captcha.Email.Valid {
		cacheKey = vars.GetEmailCapchaCacheKey(captcha.Email.String)
	} else {
		cacheKey = vars.GetPhonenumberCapchaCacheKey(captcha.PhoneNumber.String)
	}

	ret, err := m.Insert(context.Background(), &captcha)
	if err != nil {
		return nil, errors.Wrapf(err, "[InsertCaptchaToDbAndCache] insert captcha to db failed, catpcha:%+v", captcha)
	}

	if err := m.SetCacheWithExpire(cacheKey, captcha, vars.CacheExpireIn5m); err != nil {
		logx.Errorf("[InsertCaptchaToDbAndCache] set captcha to cache failed, err: %+v", err)
	}

	return ret, nil
}
