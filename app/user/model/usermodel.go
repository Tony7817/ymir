package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ymir.com/pkg/id"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		InsertIntoUserAndUserGoogle(ctx context.Context, user *User, userGoogle *UserGoogle) (int64, error)
		InsertIntoUserAndUserLocal(ctx context.Context, user *User, userLocal *UserLocal) (int64, error)
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c, opts...),
	}
}

func (m *customUserModel) InsertIntoUserAndUserGoogle(ctx context.Context, user *User, userGoogle *UserGoogle) (int64, error) {
	var userId int64
	err := m.CachedConn.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		ret, err := s.ExecCtx(ctx, "insert into `user` (id, username, email, phone_number, avatar_url, type) values (?,?,?,?,?,?)", id.SF.GenerateID(), user.Username, user.Email, user.PhoneNumber, user.AvatarUrl, UserTypeGoogle)
		if err != nil {
			return err
		}

		lastInsertId, err := ret.LastInsertId()
		if err != nil {
			return err
		}

		_, err = s.ExecCtx(ctx, "insert into user_google (id, user_id, google_user_id) values (?,?,?)", id.SF.GenerateID(), lastInsertId, userGoogle.GoogleUserId)
		if err != nil {
			return err
		}

		userId = lastInsertId

		return nil
	})
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (m *customUserModel) InsertIntoUserAndUserLocal(ctx context.Context, user *User, userLocal *UserLocal) (int64, error) {
	var userId int64
	err := m.CachedConn.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		ret, err := s.ExecCtx(ctx, "insert into `user` (id, username, email, phone_number, avatar_url, type) values (?,?,?,?,?,?)", id.SF.GenerateID(), user.Username, user.Email, user.PhoneNumber, user.AvatarUrl, UserTypeLocal)
		if err != nil {
			return err
		}

		lastInsertId, err := ret.LastInsertId()
		if err != nil {
			return err
		}

		_, err = s.ExecCtx(ctx, "insert into user_local (id, user_id, password_hash, is_activated) values (?,?,?,?)", id.SF.GenerateID(), lastInsertId, userLocal.PasswordHash, userLocal.IsActivated)
		if err != nil {
			return err
		}

		userId = lastInsertId
		return nil
	})
	if err != nil {
		return 0, err
	}

	return userId, nil
}
