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
		FindOneNoCache(ctx context.Context, id int64) (*Captcha, error)
		FindCaptchaByEmail(email string) (*Captcha, error)
		FindCaptchaByPhonenumber(phonenumber string) (*Captcha, error)
		InsertCaptchaToDbAndCache(captcha Captcha) (sql.Result, error)
		softDelete(id int64, key string) error
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

	// cache
	err := m.CachedConn.GetCache(cacheKey, &captcha)
	if errors.Is(err, nil) {
		return &captcha, nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		logx.Errorf("[FindCaptchaByEmail] get captcha from cache failed, err: %+v", err)
	}

	// db
	err = m.QueryRowNoCache(&captcha, "select * from captcha where email = ? and is_delete = 0 and created_at > DATE_SUB(now(), INTERVAL 5 MINUTE) order by created_at desc limit 1", email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	// cache
	err = m.CachedConn.SetCacheWithExpire(cacheKey, captcha, vars.CacheExpireIn5m)
	if err != nil {
		logx.Errorf("[FindCaptchaByEmail] set captcha to cache failed, err: %+v", err)
	}

	return &captcha, nil
}

func (m *customCaptchaModel) FindCaptchaByPhonenumber(phonenumber string) (*Captcha, error) {
	var cacheKey = vars.GetPhonenumberCapchaCacheKey(phonenumber)
	var captcha Captcha

	// cache
	err := m.CachedConn.GetCache(cacheKey, &captcha)
	if err == nil {
		return &captcha, nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		logx.Errorf("[FindCaptchaByPhonenumber] get captcha from cache failed, err: %+v", err)
	}

	// db
	err = m.QueryRowNoCache(&captcha, "select * from captcha where phone_number = ? and is_delete = 0 and created_at > DATE_SUB(now(), INTERVAL 5 MINUTE) order by created_at desc limit 1", phonenumber)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	// cache
	if err = m.CachedConn.SetCacheWithExpire(cacheKey, captcha, vars.CacheExpireIn5m); err != nil {
		logx.Errorf("[FindCaptchaByPhonenumber] set captcha to cache failed, err: %+v", err)
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

func (m *customCaptchaModel) FindOneNoCache(ctx context.Context, id int64) (*Captcha, error) {
	var captcha Captcha
	if err := m.QueryRowNoCacheCtx(ctx, &captcha, "select * from captcha where id = ?", id); err != nil {
		return nil, err
	}

	return &captcha, nil
}

func (m *customCaptchaModel) softDelete(id int64, key string) error {
	_, err := m.ExecNoCache("update captcha set is_delete = 1 where id = ?", id)
	if err != nil {
		return errors.Wrapf(err, "[FalseDelete] update captcha failed, id: %+v", id)
	}

	err = m.CachedConn.DelCache(key)
	if err != nil {
		logx.Errorf("[FalseDelete] delete captcha from cache failed, err: %+v", err)
	}

	return nil
}
