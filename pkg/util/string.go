package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"ymir.com/pkg/xerr"
)

func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

func MaskEmail(s string) (string, error) {
	if IsEmailValid(s) {
		s = strings.Split(s, "@")[0]
	}
	if s == "" {
		return "", errors.Wrapf(xerr.NewErrCode(xerr.ErrorReuqestParam), "mask empty string is not supported")
	}
	if len(s) < 5 {
		return string(s[0]) + "*****" + s[len(s)-1:], nil
	}

	return string(s[0]) + "*****" + s[len(s)-1:], nil
}

func GenerateCpatcha() (string, error) {
	const codeLength = 6
	const maxDigit = 10 // Digits are 0-9

	code := ""
	for i := 0; i < codeLength; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(maxDigit))
		if err != nil {
			return "", fmt.Errorf("failed to generate secure code: %w", err)
		}
		code += num.String()
	}
	return code, nil
}

// Save captcha in redis as the format of "captcha:createdAt"
func ParseRedisCaptcha(value string) (string, int64, error) {
	var vs = strings.Split(value, ":")
	if len(vs) != 2 {
		return "", 0, errors.Wrapf(xerr.NewErrCode(xerr.ErrorReuqestParam), "invalid redis captcha value")
	}

	createdAt, err := strconv.ParseInt(vs[1], 10, 64)
	if err != nil {
		return "", 0, errors.Wrapf(xerr.NewErrCode(xerr.ErrorReuqestParam), "invalid redis captcha created time")
	}

	return vs[0], createdAt, nil
}

var E164Regex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

func IsPhonenumberValid(phonenumber string) bool {
	return E164Regex.MatchString(phonenumber)
}

func MaskPhonenumber(phonenumber string) (string, error) {
	if len(phonenumber) < 4 {
		return "", xerr.NewErrCode(xerr.ErrorReuqestParam)
	}

	return "****" + phonenumber[len(phonenumber)-4:], nil
}
