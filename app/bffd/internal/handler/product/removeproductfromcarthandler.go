package product

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ymir.com/app/bffd/internal/logic/product"
	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"

	"ymir.com/pkg/result"
)

func RemoveProductFromCartHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RemoveProductFromCartRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := product.NewRemoveProductFromCartLogic(r.Context(), svcCtx)
		resp, err := l.RemoveProductFromCart(&req)
		result.HttpResult(r, w, resp, err)
	}
}
