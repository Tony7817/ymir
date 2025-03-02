package paypal

import (
	"database/sql"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
)

const PaypalGetTokenUrl = "https://api-m.sandbox.paypal.com/v1/oauth2/token"

const PaypalTokenCacheKey = "cache:paypal:access_token"

func PaypalCaptureOrderUrl(porderId string) string {
	return "https://api-m.sandbox.paypal.com/v2/checkout/orders/" + porderId + "/capture"
}

type PaypalTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type ErrorResponse struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	DebugId string `json:"debug_id"`
}

func SetErrorInNullString(errResp *ErrorResponse) sql.NullString {
	if errResp == nil {
		return sql.NullString{}
	}

	errStr, err := json.Marshal(errResp)
	if err != nil {
		logx.Errorf("marshal paypal error response failed: %v, errMsg: %+v", err, errResp)
		return sql.NullString{}
	}

	return sql.NullString{String: string(errStr), Valid: true}
}
