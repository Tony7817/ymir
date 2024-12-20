package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ymir.com/app/bffd/internal/logic/user"
	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
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
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
