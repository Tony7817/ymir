package order

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"ymir.com/app/bffd/internal/logic/order"
	"ymir.com/app/bffd/internal/svc"
	"ymir.com/app/bffd/internal/types"

	"ymir.com/pkg/result"
)

func CreatePaypalOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreatePaypalOrderRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := order.NewCreatePaypalOrderLogic(r.Context(), svcCtx)
		resp, err := l.CreatePaypalOrder(&req)
		result.HttpResult(r, w, resp, err)
	}
}
