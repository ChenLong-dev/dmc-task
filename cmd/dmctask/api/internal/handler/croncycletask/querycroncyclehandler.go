package croncycletask

import (
	"net/http"

	"dmc-task/cmd/dmctask/api/internal/logic/croncycletask"
	"dmc-task/cmd/dmctask/api/internal/svc"
	"dmc-task/cmd/dmctask/api/internal/types"
	"dmc-task/core/validators"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func QueryCronCycleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryCronCycleTaskReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// add validation logic here
		if validateErr := validators.Validate(&req); validateErr != nil {
			httpx.ErrorCtx(r.Context(), w, validateErr)
			return
		}

		l := croncycletask.NewQueryCronCycleLogic(r.Context(), svcCtx)
		resp, err := l.QueryCronCycle(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
