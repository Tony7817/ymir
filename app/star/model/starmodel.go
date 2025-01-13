package model

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ymir.com/pkg/id"
	"ymir.com/pkg/xerr"
)

var _ StarModel = (*customStarModel)(nil)

type (
	// StarModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStarModel.
	StarModel interface {
		starModel
		FindStarList(ctx context.Context, offset, limit int64) ([]Star, error)
		CountStarTotal(ctx context.Context) (int64, error)
		InsertStar(ctx context.Context, star *Star) (int64, error)
		UpdatePartial(ctx context.Context, star *StarPartial) (sql.Result, error)
	}

	customStarModel struct {
		*defaultStarModel
		sf *id.Snowflake
	}
)

type StarPartial struct {
	Id          int64            `db:"id"`
	Name        *string          `db:"name"`
	Rate        *sql.NullFloat64 `db:"rate"`
	RateCount   *int64           `db:"rate_count"`
	AvatarUrl   *string          `db:"avatar_url"`
	Description *sql.NullString  `db:"description"`
	CoverUrl    *string          `db:"cover_url"`
	PosterUrl   *string          `db:"poster_url"`
}

// NewStarModel returns a model for the database table.
func NewStarModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) StarModel {
	return &customStarModel{
		defaultStarModel: newStarModel(conn, c, opts...),
		sf:               id.NewSnowFlake(),
	}
}

func (m *customStarModel) FindStarList(ctx context.Context, offset, limit int64) ([]Star, error) {
	if limit > 20 {
		limit = 20
	}
	var stars []Star
	err := m.QueryRowsNoCacheCtx(ctx, &stars, "select * from `star` LIMIT ?,?", offset, limit)
	if err != nil {
		return nil, err
	}

	return stars, nil
}

func (m *customStarModel) CountStarTotal(ctx context.Context) (int64, error) {
	const cacheKey = "cache:star:total"
	var total int64
	if err := m.QueryRowCtx(ctx, &total, cacheKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		return conn.QueryRowCtx(ctx, &total, "SELECT COUNT(*) FROM `star`")
	}); err != nil {
		return 0, err
	}

	return total, nil
}

func (m *customStarModel) InsertStar(ctx context.Context, star *Star) (int64, error) {
	star.Id = m.sf.GenerateID()
	if _, err := m.Insert(ctx, star); err != nil {
		return 0, nil
	}

	return star.Id, nil
}

func (m *customStarModel) UpdatePartial(ctx context.Context, star *StarPartial) (sql.Result, error) {
	if star == nil {
		return nil, xerr.NewErrCode(xerr.ReuqestParamError)
	}

	var query = "update star set "
	var cond = ""
	var args []any
	if star.Name != nil {
		cond += "name = ?, "
		args = append(args, *star.Name)
	}
	if star.Rate != nil {
		cond += "rate = ?,"
		args = append(args, *star.Rate)
	}
	if star.RateCount != nil {
		cond += "rate_count = ?,"
		args = append(args, *star.RateCount)
	}
	if star.AvatarUrl != nil {
		cond += "avatar_url = ?,"
		args = append(args, *star.AvatarUrl)
	}
	if star.Description != nil {
		cond += "description = ?,"
		args = append(args, *star.Description)
	}
	if star.CoverUrl != nil {
		cond += "cover_url = ?,"
		args = append(args, *star.CoverUrl)
	}
	if star.PosterUrl != nil {
		cond += "poster_url = ?,"
		args = append(args, *star.PosterUrl)
	}
	cond = cond[:len(cond)-1]
	cond += " where id = ?"
	args = append(args, star.Id)

	if cond == "" {
		return nil, xerr.NewErrCode(xerr.ReuqestParamError)
	}

	logx.Infof("sql:%s", query+cond)

	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx, query+cond, args...)
	}, CacheKeyStarId(star.Id))
	if err != nil {
		return nil, errors.Wrap(err, "[UpdatePartial] update star failed")
	}

	return ret, err
}
