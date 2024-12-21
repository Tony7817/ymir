package user

import (
	"net/http"

	"ymir.com/app/bffd/internal/logic/user"
	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"
	"ymir.com/pkg/result"
	"ymir.com/pkg/vars"
	"ymir.com/pkg/xerr"
)

func SendSignupCaptchaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := r.Context().Value(vars.RequestContextKey).(types.SendSignupCaptchaRequest)
		if !ok {
			result.HttpResult(r, w, nil, xerr.NewErrMsg("invalid request"))
		}

		l := user.NewSendSignupCaptchaLogic(r.Context(), svcCtx)
		resp, err := l.SendSignupCaptcha(&req)
		result.HttpResult(r, w, resp, err)
	}
}
