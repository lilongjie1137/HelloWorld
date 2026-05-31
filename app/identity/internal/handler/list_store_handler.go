// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/lilongjie1137/HelloWorld/app/identity/internal/logic"
	"github.com/lilongjie1137/HelloWorld/app/identity/internal/svc"
	"github.com/lilongjie1137/HelloWorld/app/identity/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListStoreHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListStoreReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewListStoreLogic(r.Context(), svcCtx)
		resp, err := l.ListStore(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
