package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ymir.com/app/bffd/internal/logic/user"
	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/pkg/result"
)

func SendForgetPasswordCaptchaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendForgetPasswordCaptchaRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewSendForgetPasswordCaptchaLogic(r.Context(), svcCtx)
		resp, err := l.SendForgetPasswordCaptcha(&req)
		result.HttpResult(r, w, resp, err)
	}
}
