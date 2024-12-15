package util

import (
	"github.com/golang-jwt/jwt"
)

func GetJwtToken(secret string, nowDate int64, accessExpire int64, userId string) (string, error) {
	var claim = make(jwt.MapClaims)
	var token = jwt.New(jwt.SigningMethodHS256)

	claim["exp"] = nowDate + accessExpire
	claim["iat"] = nowDate
	claim["userId"] = userId
	token.Claims = claim
	return token.SignedString([]byte(secret))
}
