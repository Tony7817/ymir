package star

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ymir.com/app/bffs/internal/logic/star"
	"ymir.com/app/bffs/internal/svc"
	"ymir.com/app/bffs/internal/types"

	"ymir.com/pkg/result"
)

func StarDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.StarDetailRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := star.NewStarDetailLogic(r.Context(), svcCtx)
		resp, err := l.StarDetail(&req)
		result.HttpResult(r, w, resp, err)
	}
}
