package auth

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ymir.com/app/bffs/internal/logic/auth"
	"ymir.com/app/bffs/internal/svc"
	"ymir.com/app/bffs/internal/types"

	"ymir.com/pkg/result"
)

func SigninHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SigninRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := auth.NewSigninLogic(r.Context(), svcCtx)
		resp, err := l.Signin(&req)
		result.HttpResult(r, w, resp, err)
	}
}
