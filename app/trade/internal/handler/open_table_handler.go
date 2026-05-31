// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/lilongjie1137/HelloWorld/app/trade/internal/logic"
	"github.com/lilongjie1137/HelloWorld/app/trade/internal/svc"
	"github.com/lilongjie1137/HelloWorld/app/trade/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func OpenTableHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OpenTableReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewOpenTableLogic(r.Context(), svcCtx)
		resp, err := l.OpenTable(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
