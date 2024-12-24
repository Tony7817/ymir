package product

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ymir.com/app/bffd/internal/logic/product"
	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"

	"ymir.com/pkg/result"
)

func DecreaseProductAmountInCartHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DecreaseProductAmountInCartRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := product.NewDecreaseProductAmountInCartLogic(r.Context(), svcCtx)
		resp, err := l.DecreaseProductAmountInCart(&req)
		result.HttpResult(r, w, resp, err)
	}
}
