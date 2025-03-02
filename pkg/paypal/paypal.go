package paypal

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

func GetToken(ctx context.Context, redis *redis.Redis) (string, error) {
	tokenStr, err := redis.GetCtx(ctx, PaypalTokenCacheKey)
	if err != nil {
		return "", err
	}

	var token = &PaypalTokenResponse{}
	if tokenStr == "" {
		token, err = getPaypalToken()
		if err != nil {
			return "", err
		}
		_, err = redis.SetnxExCtx(ctx, PaypalTokenCacheKey, token.AccessToken, token.ExpiresIn-60)
		if err != nil {
			return "", err
		}
	} else {
		token.AccessToken = tokenStr
	}

	return token.AccessToken, nil
}

func getPaypalToken() (*PaypalTokenResponse, error) {
	var clientId = os.Getenv("PAYPAL_CLIENT_ID")
	var clientSecret = os.Getenv("PAYPAL_SECRET")
	if clientId == "" || clientSecret == "" {
		return nil, errors.New("paypal client id or secret is empty")
	}

	var auth = base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientSecret))
	var data = url.Values{}
	data.Set("grant_type", "client_credentials")
	var body = strings.NewReader(data.Encode())

	req, err := http.NewRequest("POST", PaypalGetTokenUrl, body)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var client = &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResp PaypalTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return nil, err
	}

	return &tokenResp, nil
}
