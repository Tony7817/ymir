package star

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ymir.com/app/bffs/internal/logic/star"
	"ymir.com/app/bffs/internal/svc"
	"ymir.com/app/bffs/internal/types"

	"ymir.com/pkg/result"
)

func CreateStarHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreteStarRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := star.NewCreateStarLogic(r.Context(), svcCtx)
		resp, err := l.CreateStar(&req)
		result.HttpResult(r, w, resp, err)
	}
}
