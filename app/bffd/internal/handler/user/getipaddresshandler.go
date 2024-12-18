package user

import (
	"net/http"

	"ymir.com/app/bffd/internal/logic/user"
	"ymir.com/app/bffd/internal/svc"

	"ymir.com/pkg/result"
)

func GetIpAddressHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetIpAddressLogic(r.Context(), svcCtx)
		resp, err := l.GetIpAddress(r)
		result.HttpResult(r, w, resp, err)
	}
}
