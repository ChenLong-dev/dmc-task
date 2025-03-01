package handler

import (
	"net/http"

	"dmc-task/cmd/frontend/biz/internal/logic"
	"dmc-task/cmd/frontend/biz/internal/svc"
	"dmc-task/cmd/frontend/biz/internal/types"
	"dmc-task/core/validators"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func BizHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// add validation logic here
		if validateErr := validators.Validate(&req); validateErr != nil {
			httpx.ErrorCtx(r.Context(), w, validateErr)
			return
		}

		l := logic.NewBizLogic(r.Context(), svcCtx)
		resp, err := l.Biz(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
