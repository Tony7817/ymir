package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ymir.com/app/bffd/internal/logic/user"
	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"

	"ymir.com/pkg/result"
)

func CheckIfUserSignedUpHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckIfUserSignedUpRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewCheckIfUserSignedUpLogic(r.Context(), svcCtx)
		resp, err := l.CheckIfUserSignedUp(&req)
		result.HttpResult(r, w, resp, err)
	}
}
