package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"ymir.com/app/bffs/internal/config"
	"ymir.com/app/user/admin/user"
	"ymir.com/app/user/admin/userclient"
	"ymir.com/pkg/vars"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type AuthMiddleware struct {
	UserRPC userclient.User
}

type OrganizerKey string

func NewAuthMiddleware(c config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		UserRPC: userclient.NewUser(zrpc.MustNewClient(c.UserAdminRPC)),
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uId, err := r.Context().Value(vars.UserIdKey).(json.Number).Int64()
		if err != nil {
			logx.Errorf("[AuthMiddleware] parse user id failed: %+v", err)
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
			return
		}

		respb, err := m.UserRPC.GetOrganizer(r.Context(), &user.GetOrganizerRequest{
			UserId: &uId,
		})
		if err != nil {
			logx.Errorf("[AuthMiddleware] GetOrganizer failed: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		const key vars.ContextKey = vars.OrganizerKey
		var ctx = context.WithValue(r.Context(), key, respb.Organizer)
		r = r.WithContext(ctx)

		next(w, r)
	}
}
