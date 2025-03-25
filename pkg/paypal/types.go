package paypal

import (
	"database/sql"
	"encoding/json"
	"os"

	"github.com/zeromicro/go-zero/core/logx"
	"ymir.com/pkg/vars"
)

const PaypalAPIPrefix = "https://api-m.paypal.com"
const PaypalAPISandboxPrefix = "https://api-m.sandbox.paypal.com"
const PaypalTokenCacheKey = "cache:paypal:access_token"

var APIPrefix = ""

func Init() {
	var mode = os.Getenv("YMIR_MODE")
	if mode == vars.ModeDev {
		APIPrefix = PaypalAPISandboxPrefix
	} else if mode == vars.ModeProd {
		APIPrefix = PaypalAPIPrefix
	} else {
		panic("invalid mode")
	}
}

func PaypalCreateOrderUrl() string {
	return APIPrefix + "/v2/checkout/orders"
}
func PaypalGetTokenUrl() string {
	return APIPrefix + "/v1/oauth2/token"
}

func PaypalCheckoutUrl() string {
	return APIPrefix + "/v2/checkout/orders"
}

func PaypalCaptureOrderUrl(porderId string) string {
	return APIPrefix + "/v2/checkout/orders/" + porderId + "/capture"
}

func PaypalShowOrderDetailUrl(orderId string) string {
	return PaypalCheckoutUrl() + "/" + orderId
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
