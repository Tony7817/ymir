package model

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		InsertIntoUserAndUserGoogle(ctx context.Context, user *User, userGoogle *UserGoogle) error
		InsertIntoUserAndUserLocal(ctx context.Context, user *User, userLocal *UserLocal) error
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

func (m *customUserModel) InsertIntoUserAndUserGoogle(ctx context.Context, user *User, userGoogle *UserGoogle) error {
	err := m.CachedConn.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		_, err := s.ExecCtx(ctx, "insert into `user` (id, username, email, phone_number, avatar_url, type) values (?,?,?,?,?,?)", user.Id, user.Username, user.Email, user.PhoneNumber, user.AvatarUrl, UserTypeGoogle)
		if err != nil {
			return err
		}

		_, err = s.ExecCtx(ctx, "insert into user_google (id, user_id, google_user_id) values (?,?,?)", userGoogle.Id, user.Id, userGoogle.GoogleUserId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *customUserModel) InsertIntoUserAndUserLocal(ctx context.Context, user *User, userLocal *UserLocal) error {
	err := m.CachedConn.TransactCtx(ctx, func(ctx context.Context, s sqlx.Session) error {
		_, err := s.ExecCtx(ctx, "insert into `user` (id, username, email, phone_number, avatar_url, type) values (?,?,?,?,?,?)", user.Id, user.Username, user.Email, user.PhoneNumber, user.AvatarUrl, UserTypeLocal)
		if err != nil {
			return err
		}

		_, err = s.ExecCtx(ctx, "insert into user_local (id, user_id, password_hash, is_activated) values (?,?,?,?)", userLocal.Id, user.Id, userLocal.PasswordHash, userLocal.IsActivated)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return errors.Wrapf(err, "insert user and user local failed")
	}

	return nil
}
