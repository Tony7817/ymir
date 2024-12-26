package result

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
	"ymir.com/pkg/xerr"
)

// http response
func HttpResult(r *http.Request, w http.ResponseWriter, resp any, err error) {
	if err == nil {
		//return success
		r := Success(resp)
		httpx.WriteJson(w, http.StatusOK, r)
	} else {
		errcode := xerr.ServerCommonError
		//default error msg
		errmsg := "Internal server error"

		causeErr := errors.Cause(err)
		if e, ok := causeErr.(*xerr.CodeError); ok {
			//custom error
			errcode = e.GetErrCode()
			errmsg = e.GetErrMsg()
		} else {
			if gstatus, ok := status.FromError(causeErr); ok {
				//grpc error by uint32 convert
				grpcCode := uint32(gstatus.Code())
				if xerr.IsCodeErr(grpcCode) {
					errcode = grpcCode
					errmsg = gstatus.Message()
				}
			}
		}
		httpx.WriteJson(w, http.StatusBadRequest, Error(errcode, errmsg))
	}
}
