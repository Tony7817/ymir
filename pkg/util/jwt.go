package util

import (
	"github.com/golang-jwt/jwt"
	"ymir.com/pkg/vars"
)

func GetJwtToken(secret string, nowDate int64, accessExpire int64, userId int64) (string, error) {
	var claim = make(jwt.MapClaims)
	var token = jwt.New(jwt.SigningMethodHS256)

	claim["exp"] = nowDate + accessExpire
	claim["iat"] = nowDate
	claim[vars.UserIdKey] = userId
	token.Claims = claim
	return token.SignedString([]byte(secret))
}
